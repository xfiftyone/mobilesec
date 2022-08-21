function hook_request(className, funcName)
{
    var func = ObjC.classes[className][funcName];
    var origImp = func.implementation;
    func.implementation = ObjC.implement(func, function  (self,sel,a1, a2, a3,a4,a5,a6){
        var new_dict = ObjC.classes.NSMutableDictionary.alloc().init();
        var old_dict = new ObjC.Object(a5);
        console.log("\nold_dict--->：\n"+old_dict.toString());
        try {
          var enumerator = old_dict.keyEnumerator();
          var key;
          while ((key = enumerator.nextObject()) !== null) {
            var value = old_dict.objectForKey_(key);
          if ("xxx" == key) {
              new_dict.setObject_forKey_('2022',"xxx");
            }else{
              new_dict.setObject_forKey_(value,key);
            }
          }
          console.log("new_dict--->：\n"+new_dict.toString());
          return origImp(self,sel,a1, a2, a3,a4,new_dict,a6);
        } catch (error) {
          console.log(error);
          return origImp(self,sel,a1, a2, a3,a4,a5,a6);
        }

    });
}

if (ObjC.available) {
    hook_request("XXXXX","- sendRequestWithUrlString:requestType:contentType:configEncryptBlock:body:andCompletionBlock:")
    hook_response()
}
