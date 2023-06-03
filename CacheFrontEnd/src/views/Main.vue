<script setup>
import { onMounted } from 'vue';
import { callPerPeriod } from '../utils';

const { proxy } = getCurrentInstance();

function getKey(key) {
  proxy.$http
    .httpRequest({
      url: 'http://localhost:5999/cache/' + key,
      method: 'get',
    })
    .then(
      (res) => {
        let value = res.data;
        console.log('get key success', key, value);
        ElMessage.success('GET key success');
        getForm.value = value;
      },
      (reason) => {
        console.log('get key failed', reason);
        ElMessage.error('get key failed.');
        getForm.value = 'cache miss';
      }
    );
}

function putKey(key, value) {
  proxy.$http
    .httpRequest({
      url: 'http://localhost:5999/cache/' + key,
      method: 'post',
      data: value,
    })
    .then(
      (res) => {
        console.log('put key success', key, value);
        ElMessage.success('PUT key success');
      },
      (reason) => {
        console.log('put key failed', reason);
        ElMessage.error('put key failed.');
      }
    );
}

function deleteKey(key) {
  proxy.$http
    .httpRequest({
      url: 'http://localhost:5999/cache/' + key,
      method: 'delete',
    })
    .then(
      (res) => {
        console.log('delete key success', key);
        ElMessage.success('DELETE key success');
      },
      (reason) => {
        console.log('delete key failed', reason);
        ElMessage.error('delete key failed.');
      }
    );
}

const getForm = reactive({
  key: '',
  value: '',
});

const putForm = reactive({
  key: '',
  value: '',
});

const deleteForm = reactive({
  key: '',
});

const dataList = reactive({
  dataList: [[], [], []],
});

function getSingleNodeData(port, number) {
  proxy.$http
    .httpRequest({
      url: 'http://localhost:' + port + '/cache/',
      method: 'put',
    })
    .then(
      (res) => {
        // console.log('getSingleNodeData success', port, res.data);
        dataList.dataList[number] = res.data.sort();
      },
      (reason) => {
        // console.log('getSingleNodeData failed', reason);
      }
    );
}

function refreshData() {
  getSingleNodeData(5001, 0);
  getSingleNodeData(5002, 1);
  getSingleNodeData(5003, 2);
}

onMounted(() => {
  callPerPeriod(500, -1, refreshData);
});
</script>

<template>
  <div class="forms">
    <div class="op-form">
      <el-form :model="getForm" label-width="60px">
        <el-form-item label="Key">
          <el-input v-model="getForm.key" />
        </el-form-item>
        <el-form-item label="Value">
          <div>{{ getForm.value }}</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getKey(getForm.key)">GET</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="op-form">
      <el-form :model="putForm" label-width="60px">
        <el-form-item label="Key">
          <el-input v-model="putForm.key" />
        </el-form-item>
        <el-form-item label="Value">
          <el-input v-model="putForm.value" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="putKey(putForm.key, putForm.value)"
            >PUT</el-button
          >
        </el-form-item>
      </el-form>
    </div>

    <div class="op-form">
      <el-form :model="deleteForm" label-width="60px">
        <el-form-item label="Key">
          <el-input v-model="deleteForm.key" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="deleteKey(deleteForm.key)"
            >DELETE</el-button
          >
        </el-form-item>
      </el-form>
    </div>
  </div>

  <div class="data-list">
    <div class="all-data">
      <div class="data-list-title">Node 1 port: 5001</div>
      <div v-for="item in dataList.dataList[0]">
        <div>{{ item }}</div>
      </div>
    </div>
    <div class="all-data">
      <div class="data-list-title">Node 2 port: 5002</div>
      <div v-for="item in dataList.dataList[1]">
        <div>{{ item }}</div>
      </div>
    </div>
    <div class="all-data">
      <div class="data-list-title">Node 3 port: 5003</div>
      <div v-for="item in dataList.dataList[2]">
        <div>{{ item }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.forms {
  display: flex;
}

.op-form {
  margin: 20px;
  width: 200px;
}

.data-list {
  display: flex;
}

.all-data {
  margin: 30px;
  padding: 5px;
  border: 1px #ababab solid;
}

.data-list-title {
  font-size: large;
  font-weight: 900;
}
</style>
