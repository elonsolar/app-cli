<script setup>


import axios from 'axios';
import vueQr from 'vue-qr/src/packages/vue-qr.vue'
import { ref, onMounted, render } from 'vue'
import { CreatePackTask, ListPackTask } from '@/api/pack'
import { NButton, useMessage, NDataTable, NSpace, NCard } from "naive-ui";
import { defineComponent, h } from "vue";
import CreatePackTaskForm from './CreatePack'

const columns = [
  {
    title: "打包日期",
    key: "CreatedAt",
    render(row) {
      return new Date(row.CreatedAt).toLocaleString()
    }
  },
  {
    title: "项目名称",
    key: "projectName"
  },
  {
    title: "说明",
    key: "description"
  },
  {
    title: "状态",
    key: "status"
  },
  {
    title: "二维码",
    key: "qrCode",
    render(row) {
      if (row.status == "打包成功") {
        return h(vueQr, {
          text: getQRUrl(row),
          size: 100
        });
      }
    }
  },

  // {
  //   title: "下载次数",
  //   key: "downloadCount"
  // },
  {
    title: "操作",
    key: "actions",
    render(row) {
      if (row.status == "打包成功") {
        const downloadLink = `${window.location.protocol}//${window.location.host}/api/static/download/app-release-${row.ID}.apk`;
        return h('a', {
          href: downloadLink,
          target: '_blank',
          class: 'download-link',
          innerHTML: '下载'
        });
      }
    }
  }
]

const data = ref([])
const showModal = ref(false)
const http = axios.create();

defineProps({
  msg: String,
})

function addPackTask() {
  showModal.value = true

}

function getQRUrl(row) {
        return `${window.location.protocol}//${window.location.host}/api/static/download/app-release-${row.ID}.apk`;
}

onMounted(() => {
  ListPackTask().then(res => {
    let respData = res.data
    if (respData.code == 200) {
      data.value = respData.data
    }

    console.log(respData)
  })
})

function download(row) {


}

async function downloadStream(row) {
  if (row.status != "打包成功") {
    // message.error("打包未完成，请等待打包完成")
    return
  }
  const response = await http({
    method: 'post',
    url: '/api/packTask/download',
    data: { "id": row.ID },
    responseType: 'blob'
  });
  const res = response.data;
  let blob = new Blob([res], { type: 'application/vnd.android.package-archive' });
  let downloadElement = document.createElement("a");
  let href = window.URL.createObjectURL(blob);
  downloadElement.href = href;
  downloadElement.download = response.headers["filename"];
  document.body.appendChild(downloadElement);
  downloadElement.click();
  document.body.removeChild(downloadElement);
  window.URL.revokeObjectURL(href);
}


</script>

<template>
  <div>
    <n-card title="" size="large">
      <p class="text-center text-xl font-bold">游翼科技APP离线打包</p>
    </n-card>
    <n-card title="" size="large">
      <n-space class="mb-3 mt-9">
        <n-button type="info" @click="addPackTask" dashed>开始打包</n-button>
      </n-space>
      <n-data-table class="mt-1" :columns="columns" :data="data" :pagination="false" :bordered="false">
        <!-- 使用具名插槽来自定义 qrCode 列的内容 -->
        <!-- <template #body="{ row }">
          <n-table-column prop="qrCode">
            <vue-qr :text="getQRUrl(row)" :size="50"></vue-qr>
          </n-table-column>
        </template> -->
      </n-data-table>
    </n-card>

    <create-pack-task-form v-model:showModal="showModal"></create-pack-task-form>
  </div>
</template>


<style scoped>
.read-the-docs {
  color: #888;
}

.download-link {
  color: blue !important;
  /* 设置链接颜色为蓝色 */
  text-decoration: none;
  /* 可选：移除下划线 */
}
</style>
