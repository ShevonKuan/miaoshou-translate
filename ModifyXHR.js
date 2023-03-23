// ==UserScript==
// @name         Modify XHR Request URL
// @namespace    http://tampermonkey.net/
// @version      1.5
// @description  nil
// @author       ShevonKwan
// @match        https://erp.91miaoshou.com/common_collect_box/items
// @icon         https://www.google.com/s2/favicons?sz=64&domain=91miaoshou.com
// @grant        GM_xmlhttpRequest
// @license      GPLv3
// ==/UserScript==



(function () {
    var open = XMLHttpRequest.prototype.open;
    XMLHttpRequest.prototype.open = function (method, url, async, user, pass) {
        if (url === "https://erp.91miaoshou.com/api/move/common_collect_box/translateCommonBoxDetail") {
            var newUrl = "https://miaoshou-translate.vercel.app/api/index";
            var headers = arguments[5];
            headers = headers ? headers : {};
            headers["Referer"] = "https://erp.91miaoshou.com";
            arguments[5] = headers;
            open.call(this, method, newUrl, async, user, pass);
        } else {
            open.call(this, method, url, async, user, pass);
        }
    };
})();



