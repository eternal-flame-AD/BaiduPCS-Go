function getSign(uri) {
    var token = document.cookie.match(/token=(\w+)/)[1]
    
    return md5(uri + ":"+ token)
}
function makeNonce() {
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    for (var i = 0; i < 15; i++)
        text += possible.charAt(Math.floor(Math.random() * possible.length));

    return text;
}