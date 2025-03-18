<template>
  <v-container>
    <v-container>
      <v-card class="mx-auto">
        <v-card-text>
          <v-row class="text-center">
            <!-- 图 SVG 区 -->
            <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 6" class="align-self-center">
              <graph-algo-graphic ref="algoGraph"></graph-algo-graphic>
              <v-col class="text-center my-auto">
                <v-btn text color="warning" :disabled="buttonDisabled" @click="showEditor()"> 编辑图边 </v-btn>
                <v-btn text color="primary" :disabled="buttonDisabled" @click="codeSubmit"> 提交代码 </v-btn>
                <v-btn text color="success" :disabled="buttonDisabled" @click="graphNextStep"> 单步运行 </v-btn>
                <v-btn text color="success" :disabled="buttonDisabled" @click="runCode"> 全部执行 </v-btn>
                <v-btn text color="error" @click="resetStep"> 重新执行 </v-btn>
              </v-col>
            </v-col>
            <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 6">
              <!-- 按钮、代码块区 -->
              <v-select v-model="selectedAlgo" :items="algoOptions" label="选择图算法"></v-select>
              <!-- <variable-table ref="varTable"></variable-table> -->
              <code-tracer ref="code_tracer" :code="code" :input-default="inputDefault"></code-tracer>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-container>
    <v-dialog v-model="dialogVisible" max-width="450px" title="编辑图信息">
      <graph-editor :initial-links="this.links" :initial-nodes="this.nodes" @changeGraph="generateGraph"></graph-editor>
    </v-dialog>
    <v-overlay :value="loading">
      <v-progress-circular :size="70" :width="7" color="purple" indeterminate></v-progress-circular>
    </v-overlay>
  </v-container>
</template>

<script>
import graphAlgoGraphic from '@/components/visualization/GraphAlgoGraphic.vue';
import codeTracer from '@/components/visualization/CodeTracer.vue';
import { getRequest, postRequest } from '@/api/request';
import graphEditor from '@/components/visualization/GraphEditor.vue';
// import variableTable from '@/components/visualization/VariableTable.vue';

export default {
  name: 'GraphAlgoView',
  components: {
    // variableTable,
    codeTracer,
    graphAlgoGraphic,
    graphEditor
  },
  data() {
    return {
      loading: false,
      isRunning: false,
      nodes: [
        { id: 0 },
        { id: 1 },
        { id: 2 },
        { id: 3 },
        { id: 4 },
        { id: 5 },
        { id: 6 },
        { id: 7 },
        { id: 8 },
        { id: 9 }
      ],
      links: [
        { source: 0, target: 1, weight: '1' },
        { source: 0, target: 5, weight: '1' },
        { source: 1, target: 2, weight: '4' },
        { source: 1, target: 7, weight: '5' },
        { source: 2, target: 3, weight: '1' },
        { source: 2, target: 9, weight: '4' },
        { source: 3, target: 4, weight: '' },
        { source: 3, target: 8, weight: '' },
        { source: 4, target: 5, weight: '' },
        { source: 4, target: 6, weight: '' },
        { source: 6, target: 7, weight: '' },
        { source: 6, target: 8, weight: '' },
        { source: 7, target: 9, weight: '' },
        { source: 8, target: 9, weight: '' }
      ],
      curStep: 0,
      buttonDisabled: false,
      selectedAlgo: null,
      algoOptions: ['广度优先搜索', '深度优先搜索', 'kruskal算法', 'dijkstra算法', 'prim算法'],
      code: '11111111',
      inputDefault: [],
      choice: 'bfs',
      dataReady: false,
      nodeSteps: [[]],
      nodeLabelSteps: [[]],
      linkSteps: [[]],
      linkTextSteps: [[]],
      dialogVisible: false,
      isEnding: false,
      timeoutIDs: [],
      highLightSteps: ['0', '1', '2', '3', '4', '5', '6', '7', '8'],
      varSteps: [
        {
          varNames: ['a', 'b', 'c', 'd', 'e'],
          values: [1, 2, 3, 4, 5]
        },
        {
          varNames: ['a', 'b', 'c', 'd', 'e'],
          values: [11, 21, 31, 41, 51]
        },
        {
          varNames: ['str1', 'bool3', 'bool2', 'd', 'e'],
          values: ['aaaa', 'true', 'false', 4, 5]
        },
        {
          varNames: ['a', 'b', 'c', 'd', 'e'],
          values: [1, 2, 3, 4, 5]
        },
        {
          varNames: ['a', 'b', 'c', 'd', 'e'],
          values: [11, 21, 31, 41, 51]
        },
        {
          varNames: ['str1', 'bool3', 'bool2', 'd', 'e'],
          values: ['aaaa', 'true', 'false', 4, 5]
        }
      ]
    };
  },
  mounted() {
    this.$store.commit('userGuide/open', this.$route.path);
    // const stid = setTimeout(() => {
    //   this.$store.commit('userGuide/close');
    // }, 5000);
    // this.timeoutIDs.push(stid);
    this.choice = 'bfs';
    if (this.$route.params.defaultAlgo) {
      this.choice = this.$route.params.defaultAlgo;
    }
    this.loadAlgoInfo();
    // this.resetStep();
    //// console.log(this.links);
    this.$refs.code_tracer.steps = this.highLightSteps;
  },
  methods: {
    generateGraph(nodesInfo, linksInfo) {
      this.dialogVisible = false;
      this.nodes = nodesInfo;
      this.links = linksInfo;
      //// console.log(this.links);
      this.$refs.algoGraph.changeData(this.nodes, this.links);
      this.codeSubmit();
    },
    async codeSubmit() {
      let codeBlanks = [];
      this.loading = true;
      for (let i = 0; i < this.nodes.length; i++) {
        this.nodes[i].id = Number(this.nodes[i].id);
      }
      for (let i = 0; i < this.links.length; i++) {
        this.links[i].source = Number(this.links[i].source);
        this.links[i].target = Number(this.links[i].target);
        this.links[i].weight = Number.isInteger(this.links[i].weight) ? Number(this.links[i].weight) : '';
      }
      if (this.$refs.code_tracer) {
        codeBlanks = JSON.parse(JSON.stringify(this.$refs.code_tracer.codeSubmit()));
        if (
          this.$refs.algoGraph.selectNode1 !== null &&
          (this.choice === 'bfs' || this.choice === 'dfs' || this.choice === 'dijkstra' || this.choice === 'prim')
        ) {
          codeBlanks[0] = this.$refs.algoGraph.selectNode1.toString();
        }
      }
      await postRequest(`/visualizations/graph/submitCode/${this.choice}`, {
        codeBlanks: codeBlanks,
        graphInfo: this.normalizeGraphInfo(this.nodes, this.links)
      })
        .then((response) => {
          let res = response.data;
          this.nodeSteps = res['nodeColors'];
          this.linkSteps = res['linkColors'];
          this.nodeLabelSteps = res['nodeLabelColors'];
          this.linkTextSteps = res['linkTextColors'];
          this.curStep = 0;
          if (res['highlightSteps']) {
            this.highLightSteps = res['highlightSteps'];
            this.$refs.code_tracer.steps = res['highlightSteps'];
          }
          // if (res['varSteps']) {
          //   this.$refs.varTable.setSteps(res['varSteps']);
          // }
        })
        .catch((error) => {
          this.$store.dispatch('snackbar/error', error);
        });
      this.loading = false;
    },
    normalizeGraphInfo(paraNodes, paraLinks) {
      // let normalizedNodes = [];
      // let normalizedLinks = [];
      // let id2index = new Array(30).fill(-1);
      // for (let i = 0; i < paraNodes.length; i++) {
      //   normalizedNodes.push({ id: i });
      //   id2index[paraNodes[i].id] = i;
      // }
      // for (let i = 0; i < paraLinks.length; i++) {
      //   normalizedLinks.push({
      //     source: id2index[paraLinks[i].source],
      //     target: id2index[paraLinks[i].target],
      //     weight: paraLinks[i].weight
      //   });
      // }
      return { nodes: paraNodes, links: paraLinks };
    },
    changeColors() {
      // this.$refs.algoGraph.setColor(this.nodeColors, this.linkColors, this.linkTextColors, this.nodeLabelColors);
      // alert('changeColor去掉了');
    },
    loadAlgoInfo() {
      let req_url = `visualizations/graph/${this.choice}`;
      this.loading = true;
      // 发送get请求
      getRequest(req_url)
        .then((response) => {
          let res = response.data;
          if (!res) return;
          this.code = response.data['code'];
          this.inputDefault = response.data['inputDefault'];
          this.dataReady = true;
          this.$refs.code_tracer.setCode(this.code, this.inputDefault);
          if (res['nodeColors']) {
            this.nodeSteps = res['nodeColors'];
            this.linkSteps = res['linkColors'];
            this.nodeLabelSteps = res['nodeLabelColors'];
            this.linkTextSteps = res['linkTextColors'];
          }
          if (res['links']) {
            this.links = res['links'];
            this.nodes = res['nodes'];
            this.$refs.algoGraph.changeData(this.nodes, this.links);
          }
          if (res['highlightSteps']) {
            this.highLightSteps = res['highlightSteps'];
            this.$refs.code_tracer.steps = res['highlightSteps'];
          }
          // if (res['varSteps']) {
          //   this.$refs.varTable.steps = res['varSteps'];
          // }
          this.loading = false;
        })
        .catch((error) => {
          this.$store.dispatch('snackbar/error', error);
          this.loading = false;
        });
    },
    graphNextStep() {
      if (this.isEnding) return;
      //传过来的四个color数组理论上都不能为空
      let i = this.curStep;
      if (i >= this.nodeSteps.length) {
        this.$refs.algoGraph.resetColor();
        this.curStep = 0;
        return;
      }
      this.buttonDisabled = true;
      this.$refs.code_tracer.nextStep();
      // this.$refs.varTable.nextStep();

      if (this.$refs.algoGraph) {
        this.$refs.algoGraph.setColor(
          this.nodeSteps[i],
          this.linkSteps[i],
          this.nodeLabelSteps[i],
          this.linkTextSteps[i]
        );
      }
      this.curStep++;
      if (!this.isRunning) {
        const stid = setTimeout(() => {
          this.buttonDisabled = false;
        }, 200);
        this.timeoutIDs.push(stid);
      }
    },
    resetStep() {
      this.isEnding = true;
      this.$refs.code_tracer.resetStep();
      // this.$refs.varTable.resetStep();
      const stid = setTimeout(() => {
        this.$refs.algoGraph.changeData(this.nodes, this.links);
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
      if (this.curStep < this.nodeSteps.length) {
        if (!this.isEnding) {
          this.graphNextStep();
          const stid = setTimeout(this.runCode, 420);
          this.timeoutIDs.push(stid);
        }
      } else {
        const stid = setTimeout(() => {
          this.isRunning = false;
          this.buttonDisabled = false;
        }, 100);
        this.timeoutIDs.push(stid);
      }
    },
    showEditor() {
      this.dialogVisible = true;
    },
    addNode() {
      if (this.nodes.length > 21) {
        this.$store.dispatch('snackbar/error', '节点ID不能超过20');
        return;
      }
      let newID = 21;
      let idPool = Array(21).fill(true);
      for (let i = 0; i < this.nodes.length; i++) {
        idPool[this.nodes[i].id] = false;
      }
      for (let i = 0; i < 21; i++) {
        if (idPool[i]) {
          newID = i;
          break;
        }
      }

      this.nodes.push({ id: newID });
      this.$refs.algoGraph.changeData(this.nodes, this.links);
    }
  },
  beforeDestroy() {
    this.isEnding = true;
    for (let i = 0; i < this.timeoutIDs.length; i++) {
      clearTimeout(this.timeoutIDs[i]);
    }
  },
  watch: {
    async selectedAlgo(newValue) {
      let type = '';
      if (newValue === '广度优先搜索') {
        type = 'bfs';
      } else if (newValue === '深度优先搜索') {
        type = 'dfs';
      } else if (newValue === 'kruskal算法') {
        type = 'kruskal';
      } else if (newValue === 'prim算法') {
        type = 'prim';
      } else if (newValue === 'dijkstra算法') {
        type = 'dijkstra';
      }
      // 修正代码运行中切换的bug
      if (type !== this.choice) {
        this.choice = type;
        this.dataReady = false;
        this.isEnding = true;
        await this.loadAlgoInfo();
        const stid = setTimeout(() => {
          this.curStep = 0;
          this.isEnding = false;
          this.isRunning = false;
          this.buttonDisabled = false;
        }, 450);
        this.timeoutIDs.push(stid);
      }
    }
  }
};
</script>

<style scoped></style>
