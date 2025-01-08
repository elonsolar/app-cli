<template>
    <n-modal v-model:show="showModal" :mask-closable="false" preset="dialog" title="新增打包任务" positive-text="确认" style="width: 600px; height: 400px;"
        negative-text="取消" @positive-click="onPositiveClick" @negative-click="onNegativeClick">


        <n-form ref="formRef" label-placement="left" :model="formModel" label-width="150px">
            <!-- 参数输入 -->
            <n-form-item label="项目名称" path="projectName">
                <n-select v-model:value="formModel.projectName" placeholder="选择项目名称" :options="options" />
            </n-form-item>
            <n-form-item label="描述" path="description">
                <n-input v-model:value="formModel.description" placeholder="请输入描述" />
            </n-form-item>
            <n-form-item label="离线打包资源（zip)" path="file">
                <n-upload multiple directory-dnd action="/api/packTask" :data="getUploadData"  @finish="handleFinish" :max="5">
                    <n-upload-dragger>
                        <div style="margin-bottom: 12px">
                            <n-icon size="48" :depth="3">
                                <ArchiveIcon />
                            </n-icon>
                        </div>
                        <n-text style="font-size: 16px">
                            点击或者拖动文件到该区域来上传
                        </n-text>
                        <n-p depth="3" style="margin: 8px 0 0 0">

                        </n-p>
                    </n-upload-dragger>
                </n-upload>
            </n-form-item>
        </n-form>
    </n-modal>
</template>

<script setup>
import { ref, watch, reactive } from 'vue'
import { defineProps } from 'vue'
import { NModal, NUpload, NUploadDragger, NForm, NFormItem, NInput ,NSelect} from 'naive-ui'


const showModal = defineModel("showModal")


const options= reactive([
        {
          label: 'guagua',
          value: 'guagua',
        },
        {
          label: 'liveshop',
          value: 'liveshop',
        }
    ])


const formModel = reactive({
    description: '',
    projectName:'guagua'
})




// 可选：监听 localShowModal 的变化，通知父组件
// watch(localShowModal, (newVal) => {
//     emit('update:showModal', newVal)
// })

function onPositiveClick() {
    console.log('positive')
}

function onNegativeClick() {
    localShowModal.value = false // 关闭模态框
}
// 上传接口需要的额外参数
function getUploadData() {
    return {
        description: formModel.description,
        projectName: formModel.projectName
    }
}

function handleFinish(){
    showModal.value = false
}

</script>
