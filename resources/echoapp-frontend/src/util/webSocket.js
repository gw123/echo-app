var isOpen = false


async function createSocket(url) {
    var ws = new WebSocket(url);
    return new Promise(function (resolve, reject) {
        ws.onopen = function (evt) {
            console.log("Connection open ...");
            isOpen = true
            resolve(ws)
        }
        ws.onerror = function (evt) {
            console.log("Connection on error ...", evt);
            reject(false)
        }
    })
}

async function NewClient(url) {
    var _client = {
        ws: null
    }

    var ws = await createSocket(url)
    if (!ws) {
        return false
    }
    _client.ws = ws


    _client.write = function (data) {
        _client.ws.send(data)
    }

    _client.close = function () {
       _client.ws.close()
    }
    return _client
}

export async function sendMessage(api, options, data) {
    const socketManager = window.api.require('socketManager');
    var sid = 0
    var info = await createSocket(options)
    if (!info || info.sid) {
        alert(info)
        return info
    }
    sid = info.sid

    socketManager.write({
        sid: sid,
        data: data
    }, function (ret, err) {
        if (ret.status) {
            alert(JSON.stringify(ret));
        } else {
            alert(JSON.stringify(err));
        }
    })

    socketManager.closeSocket({
        sid: sid
    }, function (ret, err) {
        if (ret.status) {
            alert(JSON.stringify(ret));
        } else {
            alert(JSON.stringify(err));
        }
    })
}
