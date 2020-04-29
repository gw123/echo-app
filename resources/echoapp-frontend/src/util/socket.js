
/***
 {
    sid:            //socket的唯一标识，字符串类型
    state:          //socket状态码，见常量里面的socket状态码，数字类型
    data:           //state为接收数据时的数据，字符串类型
    host：           //udp收到数据时发送方地址
    port：           //udp收到数据时发送方端口
}

 state
 101 //创建成功
 102 //连接成功
 103 //收到数据
 201 //创建失败
 202 //连接失败
 203 //异常断开
 204 //正常断开
 205 //发生未知错误断开
 */

async function createSocket(options) {
    const socketManager = window.api.require('socketManager');
    return new Promise(function (resolve, reject) {
        socketManager.createSocket(options, (ret, err) => {
            if (ret) {
                resolve(ret)
            } else {
                reject(err)
            }
        });
    })
}

async function Client(options) {

    const socketManager = window.api.require('socketManager');
    var _client = {
        sid: 0
    }

    var info = await createSocket(options)
    if (!info || info.sid) {
        alert(info)
        return info
    }

    _client.sid = info.sid

    _client.write = function (data, needBase64) {
        socketManager.write({
            sid: _client.sid,
            data: data
        }, function (ret, err) {
            if (ret.status) {
                alert(JSON.stringify(ret));
            } else {
                alert(JSON.stringify(err));
            }
        })

    }

    _client.close = function () {
        socketManager.closeSocket({
            sid: _client.sid
        }, function (ret, err) {
            if (ret.status) {
                alert(JSON.stringify(ret));
            } else {
                alert(JSON.stringify(err));
            }
        })
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
