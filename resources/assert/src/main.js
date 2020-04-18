import Vue from 'vue'
import App from './App.vue'

//组件列表
import { Field } from 'mint-ui';
Vue.component(Field.name, Field);
import { Button } from 'mint-ui';
Vue.component(Button.name, Button);

import MintUI from 'mint-ui'
import 'mint-ui/lib/style.css'

Vue.use(MintUI)
Vue.config.productionTip = false

new Vue({
  render: h => h(App),
}).$mount('#app')
