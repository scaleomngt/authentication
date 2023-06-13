<template>
  <div class="qrcode">
    该用户<span :class="classStyle">{{ isAdult }}</span>
  </div>
</template>

<script>
import modal from '@/plugins/modal'
import { calcData, } from '../api/qrcode/index'
export default {
    data(){
        return{
            isAdult:null,
            classStyle:'',
        }
    },
    mounted(){
      modal.loading("查询中...");
        let params = {
            id:this.$route.query.id,
            address:this.$route.query.address
        }
        calcData(params).then(res =>{
            if(res.code === 0){
                if(res.data.result == '0'){
                    this.isAdult = '已经成年'
                    this.classStyle = 'adult'
                    this.$forceUpdate()
                }else{
                    this.isAdult = '未成年'
                    this.classStyle = 'noAdult'
                    this.$forceUpdate()
                }
                this.$forceUpdate()
            }
            modal.closeLoading();
        })
    },
}
</script>

<style lang="less" scoped>
.qrcode{
    width: 100vw;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 30px;
    span{
        font-weight: 700;
        padding-left: 1rem;
    }
    .adult{
        color: green;
    }
    .noAdult{
        color: red;
    }
}
</style>