<template>
  <v-container class="pa-0">
    <v-tabs v-model="tabs" centered grow>
      <v-tab>编辑节点</v-tab>
      <v-tab>编辑无向边</v-tab>
    </v-tabs>
    <v-tabs-items v-model="tabs">
      <v-tab-item>
        <v-card flat>
          <v-card-text class="text-center">
            <div v-for="(node, index) in nodes" :key="index">
              <v-row class="mx-auto">
                <v-col :cols="4" :offset="3">
                  <v-text-field v-model="node.id" label="节点ID" :rules="idRules"></v-text-field>
                </v-col>
                <v-col class="text-left mt-3">
                  <v-btn color="error" icon @click="removeNode(index)">
                    <v-icon>mdi-delete</v-icon>
                  </v-btn>
                </v-col>
              </v-row>
            </div>
            <v-row>
              <v-col>
                <v-btn color="primary" outlined @click="addNode">添加节点</v-btn>
                <v-btn color="primary" class="ml-3" outlined @click="submit()">提交修改</v-btn>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-tab-item>
      <v-tab-item>
        <v-card flat>
          <v-card-text class="text-center">
            <div v-for="(link, index) in links" :key="index">
              <v-row>
                <v-col cols="3" offset="1">
                  <v-select v-model="link.source" :items="nodeOptions" :rules="notEmptyRule" label="起点ID"></v-select>
                </v-col>
                <v-col cols="3">
                  <v-select v-model="link.target" :items="nodeOptions" :rules="notEmptyRule" label="终点ID"></v-select>
                </v-col>
                <v-col cols="3">
                  <v-text-field v-model="link.weight" label="边权" :rules="weightRules"></v-text-field>
                </v-col>
                <v-col class="text-left mt-3" cols="1">
                  <v-btn color="error" icon @click="removeLink(index)">
                    <v-icon>mdi-delete</v-icon>
                  </v-btn>
                </v-col>
              </v-row>
            </div>
            <v-row>
              <v-col>
                <v-btn color="primary" outlined @click="addLink">添加边</v-btn>
                <v-btn color="primary" class="ml-3" outlined @click="submit()">提交修改</v-btn>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-tab-item>
    </v-tabs-items>
  </v-container>
</template>

<script>
export default {
  props: {
    initialNodes: {
      type: Array,
      required: true
    },
    initialLinks: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      buttonDisabled: false,
      nodes: [], //JSON.parse(JSON.stringify(this.initialNodes)),
      links: [], // JSON.parse(JSON.stringify(this.initialLinks)),
      tabs: null,
      idRules: [
        (value) => value !== '' || '不能为空',
        (value) => (Number.isInteger(Number(value)) && value >= 0 && value <= 20) || '节点 ID 必须为 0到 20 之间的整数',
        (value) => this.checkUnique(value) || '节点ID必须唯一'
      ],
      notEmptyRule: [(value) => value !== '' || '不能为空'],
      weightRules: [
        (value) => value !== '' || '不能为空',
        (value) =>
          (Number.isInteger(Number(value)) && value >= -1000 && value <= 1000) || '边权必须为 -1000 到 1000 之间的整数'
      ]
    };
  },
  mounted() {
    this.nodes = JSON.parse(JSON.stringify(this.initialNodes));
    this.links = JSON.parse(JSON.stringify(this.initialLinks));
    for (let i = 0; i < this.nodes.length; i++) {
      this.nodes[i].id = Number(this.nodes[i].id);
    }
    for (let i = 0; i < this.links.length; i++) {
      this.links[i].source = Number(this.links[i].source);
      this.links[i].target = Number(this.links[i].target);
      // this.links[i].weight = Number(this.links[i].weight);
    }
  },
  methods: {
    transStrToNum() {
      for (let i = 0; i < this.nodes.length; i++) {
        this.nodes[i].id = this.nodes[i].id === '' ? '' : Number(this.nodes[i].id);
      }
      for (let i = 0; i < this.links.length; i++) {
        this.links[i].source = this.links[i].source === '' ? '' : Number(this.links[i].source);
        this.links[i].target = this.links[i].target === '' ? '' : Number(this.links[i].target);
        this.links[i].weight = this.links[i].weight === '' ? '' : Number(this.links[i].weight);
      }
    },
    addNode() {
      this.nodes.push({ id: '' });
    },
    removeNode(index) {
      this.transStrToNum();
      for (let i = this.links.length - 1; i >= 0; i--) {
        if (
          Number(this.links[i].source) === Number(this.nodes[index].id) ||
          Number(this.links[i].target) === Number(this.nodes[index].id)
        ) {
          this.removeLink(i);
        }
      }
      this.nodes.splice(index, 1);
    },
    addLink() {
      this.links.push({ source: '', target: '', weight: '' });
    },
    removeLink(index) {
      this.transStrToNum();
      this.links.splice(index, 1);
    },
    submit() {
      this.transStrToNum();
      if (!this.checkSubmit()) {
        this.$store.dispatch('snackbar/error', '请检查表单是否符合规则');
        // alert('请检查表单是否符合规则');
      } else if (this.idIsUnique()) {
        this.$store.dispatch('snackbar/error', '请检查NodeID是否唯一');
        // alert('请检查NodeID是否唯一');
      } else {
        this.$emit('changeGraph', this.nodes, this.links);
      }
    },
    checkUnique(value) {
      return this.nodes.filter((node) => node.id === value).length <= 1;
    },
    idIsUnique() {
      let result = false;
      this.nodes.some((node, index) => {
        result = this.nodes.slice(index + 1).some((otherNode) => Number(otherNode.id) === Number(node.id));
        return result;
      });
      return result;
    },
    checkSubmit() {
      this.transStrToNum();
      for (let i = 0; i < this.nodes.length; i++) {
        this.nodes[i].id = Number(this.nodes[i].id);
        if (!(Number.isInteger(Number(this.nodes[i].id)) && this.nodes[i].id >= 0 && this.nodes[i].id <= 20)) {
          // // console.log(this.nodes[i]);
          return false;
        }
      }
      for (let i = 0; i < this.links.length; i++) {
        let l = this.links[i];
        if (!(l.source !== '' && l.target !== '')) {
          // // console.log(l);
          return false;
        } else if (!(Number.isInteger(Number(l.weight)) && l.weight >= -1000 && l.weight <= 1000)) {
          // // console.log(l);
          return false;
        }
      }
      return true;
    }
  },
  computed: {
    nodeOptions() {
      return this.nodes.map((node) => ({
        text: node.id,
        value: node.id
      }));
    }
  }
};
</script>
