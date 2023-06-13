<template>
  <div class="home">
    <div class="form">
      <div class="form_data">
        <div class="row">
          <span class="title">身份加密</span>
        </div>
        <div class="row" v-for="(item, i) in formDate" :key="i">
          <span class="label">
            {{ item.label }}
          </span>
          <el-input v-if="item.type == 'input'" class="input" v-model="form[item.valueKey]" :placeholder="item.placeholder"></el-input>
          <span v-else-if="item.type == 'radio'" class="input">
            <el-radio v-model="form[item.valueKey]" label="男">男</el-radio>
            <el-radio v-model="form[item.valueKey]" label="女">女</el-radio>
          </span>
          <span v-else-if="item.type == 'picker'" class="input">
            <el-date-picker
              style="width:100%"
              v-model="form[item.valueKey]"
              type="date"
              :placeholder="item.placeholder"
              value-format="yyyy-MM-dd">
            </el-date-picker>
          </span>
        </div>
      </div>
      <div class="operate">
        <button class="btn"  @click="submit">提交</button>
      </div>
    </div>
    <div id="qrcode" ref="exCodeRef"></div>
  </div>
</template>

<script>
import modal from '@/plugins/modal'
import { submitData, createAccount} from '../api/home/index'
import QRCode from "qrcodejs2";
export default {
  data(){
    return{
      form:{
        gender:'男',                                                                                                                                    
        birthdate:''                                                                                                                                    
      },                                                                                                                                    
      link:'',                                                                                                                                    
      qrcode:null,                                                                                                                                    
      formDate: [                                                                                                                                    
        {                                                                                                                                    
          label: '姓名',
          valueKey: 'name',
          placeholder: '请输入姓名',
          type:'input',
        }, {
          label: '性别',
          valueKey: 'gender',
          placeholder: '请输入性别',
          type:'radio',
        }, {
          label: '国籍',
          valueKey: 'nation',
          placeholder: '请输入国籍',
          type:'input',
        }, {
          label: '出生日期',
          valueKey: 'birthdate',
          placeholder: '请输入出生日期',
          type:'picker',
        },
      ],
    }
  },
  methods:{
    submit(){
      modal.loading("加密中");
      createAccount().then(res =>{
        let params = {
          addr:res.data.address,
          birthdate:this.form.birthdate,
          name:this.form.name,
          gender:this.form.gender,
          nation:this.form.nation
        }
        submitData(params).then(res1 =>{
          if(res1.code === 0){
            this.$refs.exCodeRef.innerHTML = ""
            this.qrcodeFun('http://8.129.130.135:8888/#/qrcode?id='+res1.data.id+'&address='+res1.data.address)
          }
          modal.closeLoading();
        })
      })
    },
    qrcodeFun (val) {
      this.qrcode = new QRCode("qrcode", {
        width: 150,
        height: 150,
        text: val,
      });
    },
  }
}
</script>
<style lang="less" scoped>
.home {
  width: 100%;height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  .form {
    width: 25%;
    .form_data{
      .row{
        width: 100%;height: 6vh;
        display: flex;
        justify-content: space-around;
        align-items: center;
        .label{
          width: 20%;
          text-align: justify;
          text-align-last: justify;
        }
        .input{
          width: 70%;
        }
        .title{
          text-align: center;
          width: 100%;
          font-size: 26px;
          font-weight: 700;
        }
      }
    }
    .operate{
      width: 100%;height: 8vh;
      display: flex;
      align-items: center;
      justify-content: center;
      .btn{
        width: 30%;height: 4vh;
        color: #fff;
        border: none;
        background-color: #409EFF;
        border-radius: 10px;
        cursor: pointer;
        &:hover {
          background-color: #66B1FF;
        }
      }
    }
  }
}
</style>