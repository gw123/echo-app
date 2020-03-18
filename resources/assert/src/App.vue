<template>
  <div id="app">
   <mt-header fixed title="接龙帮助工具"></mt-header>
   <mt-navbar v-model="selected">
      <mt-tab-item id="1">输入接龙内容</mt-tab-item>
      <mt-tab-item id="2">订单列表</mt-tab-item>
      <mt-tab-item id="3">汇总结果</mt-tab-item>
    </mt-navbar>

    <!-- tab-container -->
    <mt-tab-container v-model="selected">
      <mt-tab-container-item id="1">
        <mt-field placeholder="请输入接龙内容" type="textarea" rows="15" v-model="content"></mt-field>
        <mt-button type="default" style="width: 96%" @click="parseContent">开始解析</mt-button>
      </mt-tab-container-item>

      <mt-tab-container-item id="2">
          <div class="listItem" v-for="(order,index) in lineOrders" :key="index">
              <div class="nickname">{{order.nickname}}</div>
              <div class="content">
                  <div class="raw">{{order.raw}}</div>
                  <div class="error" v-if="order.error">{{order.error}}</div>
                  <div class="goodsList" v-else>
                      <div class="goods" v-for="(goods,index) in  order.goodsList" :key="index">
                          <input v-model="goods.goodsName" @change="changeGoodsName(goods, order)" style="width: 80px" placeholder="商品名">
                          <input v-model="goods.number" style="width: 40px">
                          <input v-model="goods.kind"   style="width: 30px">
                          {{goods.goodsName}}:{{goods.number}}{{goods.kind}}
                      </div>
                  </div>
              </div>
              <div class="add" @click="addGoods(order)">添加</div>
              <div style="clear: both"></div>
          </div>
      </mt-tab-container-item>

      <mt-tab-container-item id="3">
          <div class="summaryItem" v-for="(goods,index) in summary" :key="index">
              <div class="goodsName">{{goods.goodsName}}</div>
              <div class="content">
                  <div class="kindList">
                      <div class="kind" v-for="(kind,index) in  goods.kindMap" :key="index">
                         {{kind.number}} {{kind.kind}}
                      </div>
                  </div>
                  <div style="clear: both "></div>
              </div>
          </div>
      </mt-tab-container-item>
    </mt-tab-container>
  </div>
</template>

<script>
//import HelloWorld from './components/HelloWorld.vue'
import axios from 'axios';
import { Indicator } from 'mint-ui';

export default {
  name: 'App',
  components: {
    //HelloWorld
  },
  data(){
    return {
      content: undefined,
      selected: "1",
      lineOrders: [],
      summary: []
    }
  },
  created(){
    this.selected = "1"
  },
  //deep: true
  watch:{
      lineOrders: {
          handler (orders , oldVal) {
            var summary = []
            for(var key in orders) {
                var order = orders[key]
                if(order.goodsList instanceof Array) {
                    order.goodsList.forEach((goods) =>{
                        //console.log('order', goods.goodsName,goods.kind, goods.number)
                        if(summary[goods.goodsName]) {
                            if(summary[goods.goodsName].kindMap[goods.kind]) {
                                summary[goods.goodsName].kindMap[goods.kind].number += goods.number
                            }else {
                                summary[goods.goodsName].kindMap[goods.kind] = { kind: goods.kind, number: goods.number}
                            }
                        }else{
                            summary[goods.goodsName]  = { goodsName:goods.goodsName, kindMap:{}}
                            summary[goods.goodsName].kindMap[goods.kind] = { kind: goods.kind, number: goods.number}
                        }
                    })
                }
            }
            this.summary = Object.assign({}, summary)
          },
          deep: true
     }
  },
  methods:{
   parseContent() {
     Indicator.open({
        text: '解析中...',
        spinnerType: 'fading-circle'
     });
     console.log(this.content)
     var url = '/fenci/parseContent'
     axios.post(url, {words: this.content}).then( res => {
       Indicator.close()
       console.log(res)
       if(res.data.code === 0) {
          var data = res.data.data
          this.lineOrders = Object.assign({}, data.lineOrders)
          this.summary = data.summary
       }else{
           alert(res.data.msg)
       }
     }).catch(()=>{
         Indicator.close()
     })
   },
   changeGoodsName(goods, order){
       var hasGoodsName = order.goodsList.some((item,index)=>{
          console.log("goodsName: ",index, goods.goodsName, item.goodsName)
          return item.goodsName === goods.goodsName && index !== order.goodsList.length-1
       })
       if(hasGoodsName){
           goods.goodsName = ""
           alert("同名商品已经存在")
       }
   },
   addGoods(order) {
      var hasEmpty = order.goodsList.some(function (item) {
          return item.goodsName == "" || item.number == 0
      })

      if(!hasEmpty) {
          order.goodsList.push( {
              goodsName: '',
              number: 0,
              kind: '个'
          })
      }else {
          alert("当前有未完成的编辑")
      }
   }
 }

}
</script>

<style scoped lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 55px;
}
.listItem{
    padding: 10px 0;
    text-align: left;
    .add{
        font-weight: 600;
        color: #b93d1c;
    }
    .nickname{
        margin-bottom: 10px;
        text-align: left;
        font-weight: 600;
    }
    .content{
        .error{
            color: red;
        }
        .raw{
            color: #42b983;
            margin-bottom: 8px;
        }

        .goodsList{
            margin-bottom: 5px;
            .goods{
             min-width: 30%;
             margin-right: 3% ;
             float: left;
            }
        }
    }
}
.summaryItem{
   text-align: left;
   margin-bottom: 10px;
  .goodsName{
      font-weight: 600;
      min-width: 30%;
      float: left;
  }
  .kindList{
     .kind{
        float: left;
        min-width: 50px;
        margin-right: 15px;
     }
  }
}
</style>
