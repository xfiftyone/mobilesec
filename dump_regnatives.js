// dump_regnatives.js
// 运行: frida -U -f com.example.app -l dump_regnatives.js --no-pause
'use strict';

const REG_NATIVES_SLOT = 215;      // JNINativeInterface 下标
const ptrSize = Process.pointerSize;

// ---------- 1. 计算 RegisterNatives 指针 ----------
function getRegisterNativesPtr() {
	const env = Java.vm.getEnv();
	const funcTab = env.handle.readPointer(); // JNINativeInterface*
	return funcTab.add(REG_NATIVES_SLOT * ptrSize).readPointer();
}

// ---------- 2. 主 Hook ----------
function hookRegisterNatives() {
	const target = getRegisterNativesPtr();
	console.log('[*] RegisterNatives @ ' + target);

	Interceptor.attach(target, {
		onEnter(args) {
			this.methods = args[2];
			this.nMethods = args[3].toInt32();
			this.clazz = args[1];
		},
		onLeave() {
			const env = Java.vm.getEnv();
			const className = env.getClassName(this.clazz); // com/example/Foo

			for (let i = 0; i < this.nMethods; ++i) {
				const base = this.methods.add(i * ptrSize * 3);
				const name = base.readPointer().readCString();
				const sig = base.add(ptrSize).readPointer().readCString();
				const fnPtr = base.add(ptrSize * 2).readPointer();

				const mod = Process.findModuleByAddress(fnPtr);
				const soName = mod ? mod.name : 'anon';
				const offset = mod ? fnPtr.sub(mod.base) : ptr('0');

				console.log(
					`[+] ${className}->${name}${sig}\n` +
					`    ↪ ${soName} + 0x${offset.toString(16)} (abs ${fnPtr})`
				);
			}
		}
	});
}

// ---------- 3. 可选：dlopen 之后再扫一次 ----------
function hookDlopen() {
	const dl = Module.findExportByName(null, 'android_dlopen_ext') ||
		Module.findExportByName(null, '__loader_dlopen') ||
		Module.findExportByName(null, 'dlopen');
	if (!dl) return;

	Interceptor.attach(dl, {
		onLeave(ret) {
			if (ret.isNull()) return;
			// 新 so 进来后，立即刷新函数表指针，确保 hook 仍在
			const newPtr = getRegisterNativesPtr();
			if (!Interceptor.replace || Interceptor.revert == null) return; // Frida <16
			// nothing else needed: 已经 attach 过的地址依旧有效
		}
	});
}

// ---------- 4. 入口 ----------
Java.perform(() => {
	hookRegisterNatives();
	hookDlopen();
});