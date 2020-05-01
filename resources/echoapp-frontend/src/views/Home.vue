<style scoped lang="scss">
    a-button {
        width: 100%;
    }

    .debug-from {
        padding: 15px 10px;
    }
</style>
<template>
    <div class="home">
        <a-row>
            <a-col span="8">
                <a-button @click="getInfo">获取设备编号</a-button>
            </a-col>
            <a-col span="8">
                <a-button @click="getInfo" icon="search">搜索打印机</a-button>
            </a-col>
            <a-col span="8">
                <a-button @click="connectWsServer" icon="search">远程协助</a-button>
            </a-col>
        </a-row>
        <div class="debug-from">
            <a-form-model ref="ruleForm" :model="ruleForm" :rules="rules">
                <a-form-model-item has-feedback label="IP地址" prop="address">
                    <a-input v-model="ruleForm.address" autocomplete="off"/>
                </a-form-model-item>
                <a-form-model-item has-feedback label="端口" prop="port">
                    <a-input v-model="ruleForm.port" type="number" autocomplete="off"/>
                </a-form-model-item>
                <a-form-model-item has-feedback label="发送内容" prop="content">
                    <a-input v-model="ruleForm.content" type="textarea " autocomplete="off"/>
                </a-form-model-item>
                <a-form-model-item :wrapper-col="{ span: 10, offset: 0 }">
                    <a-button type="primary" @click="sendTcpContent()">发送</a-button>
                    <a-button style="margin-left: 10px" @click="resetForm('ruleForm')">重置</a-button>
                    <a-button style="margin-left: 10px" @click="tcp_logs=[]">清空日志</a-button>
                </a-form-model-item>
            </a-form-model>
            <div v-for="(log,index) in tcp_logs" :key="index">
                <a-row>
                    <a-col :span="3">
                        {{log.pos}}
                    </a-col>
                    <a-col :span="21">
                        {{log.msg}}
                    </a-col>
                </a-row>
            </div>
        </div>
    </div>
</template>

<script>

    export default {
        name: 'Home',
        components: {},
        data() {
            return {
                rules: {},
                ruleForm: {
                    port: '9100',
                    address: '192.168.1.1',
                    content: 'this is test content\n\n',
                },
                tcp_logs: [],
                ws: null
            }
        },
        created() {

        },
        methods: {
            sendTcpContent(data) {
                this.tcp_logs.push({"pos": "发送消息", "msg": this.ruleForm.content})
                var address = this.ruleForm.address + ":" + this.ruleForm.port
                //alert("v3:" + address)
                var socketManager = window.api.require('socketManager');
                socketManager.createSocket({
                    host: this.ruleForm.address,
                    port: this.ruleForm.port
                }, (ret, err) => {
                    if (!ret) {
                        this.tcp_logs.push({"pos": "创建失败", "msg": err})
                        return
                    }
                    switch (ret.state) {
                        case 101:
                            this.tcp_logs.push({"pos": "创建成功", "msg": err + JSON.stringify(ret)})
                            break
                        case 201:
                            this.tcp_logs.push({"pos": "创建失败", "msg": err + JSON.stringify(ret)})
                            break
                        case 202:
                            this.tcp_logs.push({"pos": "连接失败", "msg": err + JSON.stringify(ret)})
                            break
                        case 204:
                            this.tcp_logs.push({"pos": "正常断开", "msg": err + JSON.stringify(ret)})
                            break
                        case 102:
                            this.tcp_logs.push({"pos": "连接成功", "msg": err + JSON.stringify(ret)})
                            break
                        default:
                            this.tcp_logs.push({"pos": "未知状态", "msg": err + JSON.stringify(ret)})

                    }

                    if (ret.state != 102) {
                        return;
                    }

                    socketManager.write({
                        sid: ret.sid,
                        data: this.ruleForm.content,
                        base64: false
                    }, (ret1, err1) => {
                        this.tcp_logs.push({"pos": "发送", "msg": err1 + JSON.stringify(ret1)})
                    })

                    socketManager.closeSocket({
                        sid: ret.sid
                    }, (ret2, err2) => {
                        this.tcp_logs.push({"pos": "关闭", "msg": err2 + JSON.stringify(ret2)})
                    })

                });
            },
            resetForm(data) {
                this.ruleForm = {
                    port: '9100',
                    address: '192.168.1.1',
                    content: '',
                }
            },
            getInfo() {
                //alert(window.api.deviceId)
            },
            createWsClient() {
                var url = 'ws://192.168.30.127:8082/gapi/createWsClient?token=MTAwMDAwMDAx2334016c6181afc8a5cad07e6a3b35a3'
                this.ws = new WebSocket(url);
            },
            sendWsLogEvent(content) {
                var event = {content}
                this.sendEvent('log', event)
            },
            sendPingEvent() {
                var event = {}
                this.sendEvent('ping', event)
            },
            sendEvent(event_type, event) {
                var deviceId = window.api ? window.api.deviceId : ''
                event.source = 'android_client' + ":" + deviceId
                event.event_id = new Date().getMilliseconds()
                event.event_type = event_type
                var msgStr = event_type + "#" + JSON.stringify(event)
                this.ws.send(msgStr);
            },
            closeWsClient() {
                this.ws.close()
                this.ws = null
            },
            connectWsServer() {
                if (this.ws) {
                    this.ws.close()
                }
                var url = 'ws://192.168.30.127:8082/createWsClient?token=MTAwMDAwMDAx2334016c6181afc8a5cad07e6a3b35a3'
                this.createWsClient(url)
                this.ws.onopen = (evt) => {
                    console.log("Connection open ...");
                    setInterval(() => {
                        this.sendPingEvent()
                        for (; this.tcp_logs.length > 0;) {
                            var item = this.tcp_logs.pop()
                            this.sendWsLogEvent(item.msg)
                        }
                    }, 2000)
                };

                this.ws.onmessage = (evt) => {
                    console.log("Received Message: " + evt.data);
                };

                this.ws.onclose = (evt) => {
                    setTimeout(() => {
                        this.createWsClient()
                    }, 10000)
                };

            }
        },

    }
</script>
