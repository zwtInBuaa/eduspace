<template>
  <v-container style="margin-top: 20px">
    <v-card class="mx-auto">
      <v-card-text>
        <v-row class="text-center">
          <!-- 左侧 -->
          <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 6">
            <code-tracer :input-default="[]" code="" ref="code_tracer"></code-tracer>
            <v-col class="text-center my-auto">
              <v-btn text color="warning" :disabled="buttonDisabled" @click="() => (showEditor = true)">编辑代码</v-btn>
              <v-btn text color="primary" :disabled="buttonDisabled" @click="submitCode">提交代码</v-btn>
              <v-btn text color="success" :disabled="buttonDisabled" @click="nextStep">单步运行</v-btn>
              <v-btn text color="success" :disabled="buttonDisabled" @click="runCode">全部执行</v-btn>
              <v-btn text color="error" @click="resetStep">重新执行</v-btn>
            </v-col>
          </v-col>
          <!-- 右侧 -->
          <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 6">
            <v-card flat class="mb-1">
              <v-card-title :class="captionStyle + ' pt-1'">程序输出</v-card-title>
              <v-card-text class="mx-auto">
                <v-row>
                  <v-col class="pre">{{ steps[curStep]['stdout'] }}</v-col>
                </v-row>
              </v-card-text>
            </v-card>
            <v-divider />
            <v-card flat class="my-1">
              <v-simple-table>
                <template v-slot:default>
                  <tbody>
                    <tr>
                      <th :class="captionStyle">当前事件</th>
                      <td>{{ steps[curStep]['event'] }}</td>
                    </tr>
                    <tr>
                      <th :class="captionStyle">当前函数</th>
                      <td>{{ steps[curStep]['func_name'] }}</td>
                    </tr>
                  </tbody>
                </template>
              </v-simple-table>
            </v-card>
            <v-divider />
            <v-card flat class="my-1 text-left">
              <v-tabs v-model="tabs" centered grow height="38">
                <v-tab>局部变量</v-tab>
                <v-tab>全局变量</v-tab>
              </v-tabs>
              <v-tabs-items v-model="tabs">
                <!-- 局部变量标签页 -->
                <v-tab-item>
                  <div style="overflow-y: auto">
                    <v-simple-table dense fixed-header style="min-height: 100px; max-height: 150px">
                      <tbody class="text-left">
                        <tr v-for="item in steps[curStep]['local_vars']" :key="item.key">
                          <td>{{ item.key }}</td>
                          <td>{{ item.value }}</td>
                        </tr>
                      </tbody>
                    </v-simple-table>
                  </div>
                </v-tab-item>
                <!-- 全局变量标签页 -->
                <v-tab-item>
                  <div style="overflow-y: auto">
                    <v-simple-table dense fixed-header style="min-height: 100px; max-height: 150px">
                      <template v-slot:default>
                        <tbody class="text-left">
                          <tr v-for="item in steps[curStep]['global_vars']" :key="item.key">
                            <td>{{ item.key }}</td>
                            <td>{{ item.value }}</td>
                          </tr>
                        </tbody>
                      </template>
                    </v-simple-table>
                  </div>
                </v-tab-item>
              </v-tabs-items>
            </v-card>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
    <div class="dialog-container" v-if="showEditor">
      <!--        不用v-show是因为有bug会影响codemirror-->
      <v-card width="80%" :loading="isLoading">
        <template slot="progress">
          <v-progress-linear color="deep-purple" height="10" indeterminate></v-progress-linear>
        </template>
        <v-card-title>输入 Python 代码</v-card-title>
        <v-card-text class="py-0">
          <codemirror ref="textarea" v-model="code" :options="options"></codemirror>
        </v-card-text>
        <v-card-actions class="my-0 pa-0">
          <v-col class="text-center">
            <v-btn color="primary" text @click="submitCode" :loading="isLoading">提交代码</v-btn>
            <v-btn color="error" text @click="() => (showEditor = false)" :disabled="isLoading">关闭</v-btn>
          </v-col>
        </v-card-actions>
      </v-card>
    </div>
    <v-overlay :value="isLoading">
      <v-progress-circular :size="70" :width="7" color="purple" indeterminate></v-progress-circular>
    </v-overlay>
  </v-container>
</template>

<script>
import CodeTracer from '@/components/visualization/CodeTracer.vue';
import { codemirror } from 'vue-codemirror';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/python/python';
import 'codemirror/theme/base16-light.css';
import { postRequest } from '@/api/request';
import MarkdownDisplay from '@/components/MarkdownDisplay.vue';
export default {
  name: 'CustomVisualizationView',
  // eslint-disable-next-line vue/no-unused-components
  components: { CodeTracer, codemirror, MarkdownDisplay },
  data() {
    return {
      algoId: -1,
      showEditor: false,
      code: '#在此写入代码',
      captionStyle: 'font-weight-medium text-subtitle-1',
      tabs: null,
      options: {
        tabSize: 4, // 制表符的宽度
        indentWithTabs: false, // 取消indentWithTabs
        smartIndent: true,
        lineNumbers: true, //打开会有遮挡bug
        theme: 'base16-light',
        matchBrackets: true,
        viewportMargin: Infinity, //处理高度自适应时搭配使用
        mode: 'python',
        extraKeys: {
          Tab: function (cm) {
            const spaces = Array(cm.getOption('indentUnit') + 1).join('  ');
            cm.replaceSelection(spaces);
          }
        }
      },
      isEnding: false,
      timeoutIDs: [],
      curStep: 0,
      buttonDisabled: false,
      isRunning: false,
      steps: [{ event: 'step_line', func_name: '<module>', stdout: '', global_vars: [], local_vars: [] }],
      isLoading: false
    };
  },
  mounted() {
    this.$store.commit('userGuide/open', this.$route.path);
    // const stid = setTimeout(() => {
    //   this.$store.commit('userGuide/close');
    // }, 30000);
    // this.timeoutIDs.push(stid);
    let initData = JSON.parse(
      '{"steps": [{"event": "step_line", "func_name": "<module>", "stdout": "", "global_vars": [], "local_vars": []}, {"event": "step_line", "func_name": "<module>", "stdout": "", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": []}, {"event": "call", "func_name": "hanoi", "stdout": "", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 4}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 4}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 4}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "call", "func_name": "hanoi", "stdout": "", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 3}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 3}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 3}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "call", "func_name": "hanoi", "stdout": "", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "call", "func_name": "hanoi", "stdout": "", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}, {"key": "__return__", "value": null}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "call", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}, {"key": "__return__", "value": null}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}, {"key": "__return__", "value": null}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 3}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 3}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "call", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "C"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "C"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "C"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "A"}]}, {"event": "call", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "C"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "C"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "C"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "C"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "B"}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "C"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "B"}, {"key": "__return__", "value": null}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "C"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "C"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "A"}]}, {"event": "call", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}, {"key": "__return__", "value": null}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "C"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "A"}, {"key": "__return__", "value": null}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 3}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}, {"key": "__return__", "value": null}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 4}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 4}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "call", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 3}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 3}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 3}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "call", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "B"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "B"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "B"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "C"}]}, {"event": "call", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}, {"key": "__return__", "value": null}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "B"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "B"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "C"}]}, {"event": "call", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "C"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "C"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "C"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "C"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "B"}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "C"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "B"}, {"key": "__return__", "value": null}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "B"}, {"key": "target", "value": "A"}, {"key": "auxiliary", "value": "C"}, {"key": "__return__", "value": null}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 3}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 3}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "call", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "call", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "A"}, {"key": "target", "value": "B"}, {"key": "auxiliary", "value": "C"}, {"key": "__return__", "value": null}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\nMove disk 2 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}]}, {"event": "call", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\nMove disk 2 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\nMove disk 2 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\nMove disk 2 from A to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "step_line", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}, {"key": "__return__", "value": null}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 2}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}, {"key": "__return__", "value": null}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 3}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}, {"key": "__return__", "value": null}]}, {"event": "return", "func_name": "hanoi", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 4}, {"key": "source", "value": "A"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "B"}, {"key": "__return__", "value": null}]}, {"event": "return", "func_name": "<module>", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}, {"key": "__return__", "value": null}]}, {"event": "return", "func_name": "<module>", "stdout": "Move disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\nMove disk 3 from A to B\\nMove disk 1 from C to A\\nMove disk 2 from C to B\\nMove disk 1 from A to B\\nMove disk 4 from A to C\\nMove disk 1 from B to C\\nMove disk 2 from B to A\\nMove disk 1 from C to A\\nMove disk 3 from B to C\\nMove disk 1 from A to B\\nMove disk 2 from A to C\\nMove disk 1 from B to C\\n", "global_vars": [{"key": "hanoi", "value": "FUNCTION:hanoi(n, source, target, auxiliary)"}], "local_vars": [{"key": "n", "value": 1}, {"key": "source", "value": "B"}, {"key": "target", "value": "C"}, {"key": "auxiliary", "value": "A"}, {"key": "__return__", "value": null}]}], "highlightSteps": [[0], [14], [0], [4], [8], [0], [4], [8], [0], [4], [8], [0], [4], [5], [6], [6], [9], [10], [0], [4], [5], [6], [6], [10], [9], [10], [0], [4], [8], [0], [4], [5], [6], [6], [9], [10], [0], [4], [5], [6], [6], [10], [10], [9], [10], [0], [4], [8], [0], [4], [8], [0], [4], [5], [6], [6], [9], [10], [0], [4], [5], [6], [6], [10], [9], [10], [0], [4], [8], [0], [4], [5], [6], [6], [9], [10], [0], [4], [5], [6], [6], [10], [10], [10], [14]]}'
    );
    this.$refs.code_tracer.steps = initData['highlightSteps'];
    this.steps = initData['steps'];
    // // console.log(this.steps.length);
    this.code =
      'def hanoi(n, source, target, auxiliary):\n' +
      '    """\n' +
      '    递归函数，将 n 个盘子从 source 移动到 target 上，借助 auxiliary\n' +
      '    """\n' +
      '    if n == 1:\n' +
      '        print(f"Move disk 1 from {source} to {target}")\n' +
      '        return\n' +
      '\n' +
      '    hanoi(n - 1, source, auxiliary, target)\n' +
      '    print(f"Move disk {n} from {source} to {target}")\n' +
      '    hanoi(n - 1, auxiliary, target, source)\n' +
      '\n' +
      '\n' +
      '# 调用 hanoi 函数进行测试\n' +
      "hanoi(4, 'A', 'C', 'B')\n";
    this.$refs.code_tracer.setCodeNoEdit(this.code);
  },
  beforeDestroy() {
    for (let i = 0; i < this.timeoutIDs.length; i++) {
      clearTimeout(this.timeoutIDs[i]);
    }
  },
  methods: {
    async submitCode() {
      this.$refs.code_tracer.resetStep();
      this.curStep = 0;
      this.$refs.code_tracer.setCodeNoEdit(this.code);
      this.isLoading = true;
      await postRequest(`/visualizations/custom/submitCode`, {
        code: this.code
      })
        .then((response) => {
          let res = response.data;
          // console.log(res);
          if (res['highlightSteps']) {
            this.$refs.code_tracer.steps = res['highlightSteps'];
          }
          this.steps = res['steps'];
          this.isLoading = false;
          this.showEditor = false;
        })
        .catch((error) => {
          this.$store.dispatch('snackbar/error', error);
          this.isLoading = false;
          this.showEditor = false;
        });
    },
    // insertTab(event) {
    //   // if (event.key === 'tab') {
    //   event.preventDefault();
    //   const textarea = this.$refs.textarea.$el;
    //   const start = textarea.selectionStart;
    //   const end = textarea.selectionEnd;
    //   // 在光标位置插入四个空格
    //   const spaces = '    ';
    //   this.code = this.code.substring(0, start) + spaces + this.code.substring(end);
    //   // 更新光标位置
    //   textarea.selectionStart = textarea.selectionEnd = start + spaces.length;
    //   // console.log(this.code);
    //   // }
    // },
    nextStep() {
      if (this.isEnding) return;
      //传过来的四个color数组理论上都不能为空
      let i = this.curStep;
      // // console.log('main', this.curStep);
      this.buttonDisabled = true;
      this.$refs.code_tracer.nextStep();
      this.curStep++;
      if (i >= this.steps.length - 1) {
        this.curStep = 0;
      }
      if (!this.isRunning) {
        const stid = setTimeout(() => {
          this.buttonDisabled = false;
        }, 50);
        this.timeoutIDs.push(stid);
      }
    },
    resetStep() {
      this.isEnding = true;
      this.$refs.code_tracer.resetStep();
      const stid = setTimeout(() => {
        this.curStep = 0;
        this.isEnding = false;
        this.isRunning = false;
        this.buttonDisabled = false;
      }, 300);
      this.timeoutIDs.push(stid);
    },
    runCode() {
      this.buttonDisabled = true;
      this.isRunning = true;
      if (this.curStep < this.steps.length - 1) {
        if (!this.isEnding) {
          this.nextStep();
          const stid = setTimeout(this.runCode, 300);
          this.timeoutIDs.push(stid);
        }
      } else {
        const stid = setTimeout(() => {
          this.isRunning = false;
          this.buttonDisabled = false;
        }, 100);
        this.timeoutIDs.push(stid);
      }
    }
  }
};
</script>

<style scoped>
.list-table {
  border-collapse: collapse;
  border-radius: 5px;
  overflow: hidden;
  font-family: Arial, sans-serif;
}

.table-title th,
.table-content td {
  padding: 10px;
  text-align: center;
  border: 1px solid #fcfcff;
}

.table-title th {
  background-color: #fafaff;
  /*color: #6d8fb7;*/
}

.table-content td {
  background-color: #fcfcff;
  /*color: #506784;*/
}

.table-content th {
  background-color: #fcfcff;
  /*color: #506784;*/
}

.dialog-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}
.pre {
  font-family: Consolas, Menlo, Courier, monospace;
  border-radius: 5px;
  background-color: #f5f5f5;
  padding: 5px;
  margin: 5px;
  width: 80%;
  overflow-x: auto;
  text-align: left;
  white-space: pre;
  max-height: 130px;
  height: 130px;
}
</style>
