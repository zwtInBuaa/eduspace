<template>
  <v-container>
    <v-card class="my-3 mx-auto" :width="$vuetify.breakpoint.mobile ? '90%' : '95%'">
      <v-card-text>
        <v-row class="mx-auto ml-1">
          <v-col :cols="$vuetify.breakpoint.mobile ? 6 : 4">
            <v-select
              v-model="selectedAlgoType"
              :items="algoTypes"
              label="算法类别"
              item-value="abbreviation"
              item-text="name"
              outlined
              dense
              @change="updateAlgos"
            ></v-select>
          </v-col>
          <v-col :cols="$vuetify.breakpoint.mobile ? 6 : 4">
            <v-select
              v-model="selectedAlgoName"
              :items="algos"
              label="算法"
              item-value="abbreviation"
              item-text="name"
              outlined
              dense
            ></v-select>
          </v-col>
          <v-col class="text-center my-1">
            <v-btn class="" color="success" outlined @click="go2algo()" :disabled="isEmpty()"> 快速跳转 </v-btn>
            <v-btn class="ml-4" color="primary" outlined @click="generateQRcode()" :disabled="isEmpty()">
              生成二维码
            </v-btn>
          </v-col>
        </v-row>
        <v-row>
          <v-col v-for="item in items" :key="item.title" :cols="$vuetify.breakpoint.mobile ? 6 : 3">
            <v-card :to="item.index">
              <v-img :src="item.img"> </v-img>
              <v-card-title> {{ item.title }} </v-card-title>
            </v-card>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
    <v-dialog v-model="dialogVisible" max-width="450px">
      <img :src="imgStr" alt="二维码未加载" />
      <v-btn color="success" @click="dialogVisible = false">关闭</v-btn>
    </v-dialog>
  </v-container>
</template>

<script>
import router from '@/router';
// import { getRequest } from '@/api/request';
import axios from '@/api/intercept';
const algoRouterParams = [
  '',
  ['/binarySearch', '/binarySearch'],
  ['/sort', '/bubbleSort', '/selectionSort', '/quickSort'],
  ['/graph', '/bfs', '/dfs', '/kruskal', '/prim', '/dijkstra']
];
const algoNames = [
  '交互式学习',
  ['-查找', '-二分查找'],
  ['-排序', '-冒泡排序', '-选择排序', '-快速排序'],
  ['-图与树', '-广度优先搜索', '-深度优先搜索', '-kruskal算法', '-prim算法', '-dijkstra算法']
];
const algoRoutes = [];
const algoTypes = [];
let preFix1 = [algoRouterParams[0], algoNames[0]];
for (let i = 1; i < algoRouterParams.length; i++) {
  let preFix2 = [algoRouterParams[i][0], algoNames[i][0]];
  let namesOfType = [];
  for (let j = 1; j < algoRouterParams[i].length; j++) {
    algoRoutes.push({
      route: preFix1[0] + preFix2[0] + algoRouterParams[i][j],
      name: preFix1[1] + preFix2[1] + algoNames[i][j]
    });
    namesOfType.push({
      abbreviation: preFix1[0] + preFix2[0] + algoRouterParams[i][j],
      name: algoNames[i][j].substring(1)
    });
  }
  algoTypes.push({ abbreviation: preFix1[0] + preFix2[0], name: preFix1[1] + preFix2[1], algos: namesOfType });
}
export default {
  data() {
    return {
      // algoRoutes,
      curRoute: this.$route.fullPath,
      selectedAlgoType: null,
      selectedAlgoName: null,
      algoTypes,
      algos: [],
      imgStr: '',
      dialogVisible: false
    };
  },
  computed: {
    items() {
      return [
        {
          title: '二分查找',
          index: this.curRoute + '/binarySearch',
          img: require('../../assets/search.svg')
        },
        {
          title: '冒泡排序',
          index: this.curRoute + '/sort/bubbleSort',
          img: 'https://visualgo.net/img/png/sorting.png'
        },
        {
          title: '选择排序',
          index: this.curRoute + '/sort/selectionSort',
          img: 'https://visualgo.net/img/png/sorting.png'
        },
        {
          title: '快速排序',
          index: this.curRoute + '/sort/quickSort',
          img: 'https://visualgo.net/img/png/sorting.png'
        },
        {
          title: '广度优先搜索',
          index: this.curRoute + '/graph/bfs',
          img: 'https://visualgo.net/img/png/dfsbfs.png'
        },
        {
          title: '深度优先搜索',
          index: this.curRoute + '/graph/dfs',
          img: 'https://visualgo.net/img/png/dfsbfs.png'
        },
        {
          title: 'kruskal算法',
          index: this.curRoute + '/graph/kruskal',
          img: 'https://visualgo.net/img/png/mst.png'
        },
        {
          title: 'dijkstra算法',
          index: this.curRoute + '/graph/dijkstra',
          img: 'https://visualgo.net/img/png/sssp.png'
        },
        {
          title: 'prim算法',
          index: this.curRoute + '/graph/prim',
          img: 'https://visualgo.net/img/png/mst.png'
        },
        {
          title: '自定义算法',
          index: this.curRoute + '/custom',
          img: 'https://visualgo.net/img/png/mst.png'
        }
      ];
    }
  },
  beforeMount() {
    setInterval(this.updateTime, 1000);
  },
  mounted() {
    // // send Req
    // // console.log(algoTypes);
    // let data = [{ id: 1, name: 'custom1' }, { id: 2, name: 'custom2' }];
    // let customAlgos = [];
    // for (let i = 0; i < data.length; i++) {
    //   customAlgos.push({
    //
    //   });
    // }
    // algoTypes.push({ abbreviation: '/custom', name: '交互式学习-自定义', algos: customAlgos });
  },
  methods: {
    updateAlgos() {
      const algoType = this.algoTypes.find((s) => s.abbreviation === this.selectedAlgoType);
      this.algos = algoType ? algoType.algos : [];
      this.selectedAlgoName = null;
      // // console.log(this.algos);
    },
    go2algo() {
      if (this.selectedAlgoType !== null && this.selectedAlgoName !== null) {
        router.push({ path: this.curRoute + this.selectedAlgoName });
      }
    },
    generateQRcode() {
      let urlPara = 'http://114.116.211.180:6060/#' + this.curRoute + this.selectedAlgoName;
      axios
        .get(`/utils/getqrcode?url=${encodeURIComponent(urlPara)}`, {
          responseType: 'arraybuffer'
        })
        .then((response) => {
          const bytes = new Uint8Array(response.data);
          const base64 = btoa(String.fromCharCode.apply(null, bytes));
          this.imgStr = 'data:image/png;base64,' + base64;
          // // console.log(this.imgStr);
          this.dialogVisible = true;
          // this.imgStr =
          //   'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAEAAAAAApiSv5AAAEtElEQVR4nOzdwW4DKRAA0fUq///L2ZMvKCNouhtmVfWuduyJVUIajOHn9/cfgf17+wJ0lwHAGQCcAcAZAJwBwBkAnAHAGQCcAcAZANzP0wOfz9kLGb+TGN9/9p1F9vkz0euLPt7t6fNwBIAzADgDgDMAOAOAMwA4A4B7nAcYVa8drL4Pjl5fdJ5gdr3ReYjo9UStfr6OAHAGAGcAcAYAZwBwBgBnAHDL8wCj7Pfp2ferfr1q2es99fk6AsAZAJwBwBkAnAHAGQCcAcBtzwOc1n3ff/p3CW/hCABnAHAGAGcAcAYAZwBwBgD32nmA6t/jz54fXfc/ys4j3OIIAGcAcAYAZwBwBgBnAHAGALc9D9B9X5u975/JrtPv3vfv1LyBIwCcAcAZAJwBwBkAnAHAGQDc8jzA6f3ts6r38+8+D+DW5+sIAGcAcAYAZwBwBgBnAHAGAPd563r1UXRdfvb1o/f1/5fPceQIAGcAcAYAZwBwBgBnAHAGAPc4D5C97z19nkBU9zr/t/3u4Ol6HAHgDADOAOAMAM4A4AwAzgDgHn8XcPo+P/t6p8/py15Pdl/CKo4AcAYAZwBwBgBnAHAGAGcAcI/zANX3oaf32+/+Pr/6/6k+X2D1/3cEgDMAOAOAMwA4A4AzADgDgHvNeoCo6vUD2dePqp4X2L1eRwA4A4AzADgDgDMAOAOAMwC45fUAs/vQ7Pn82XmC2+cGRkU/z+rX/3IEgDMAOAOAMwA4A4AzADgDgNs+N7B6P/3u/QEo+/6tvv+XIwCcAcAZAJwBwBkAnAHAGQDc8u8CZt/vV/+efubWeftP7z9ef/f3/VXzDI4AcAYAZwBwBgBnAHAGAGcAcNfODcx+P377Prp7v4PqeRHXA+hPBgBnAHAGAGcAcAYAZwBwy78LGFXft57eX796nX72eqvW+a8+/uUIAGcAcAYAZwBwBgBnAHAGALc9D5CVva9/27l83ecRdO1r6AgAZwBwBgBnAHAGAGcAcAYAt31u4Kj7PP3Z60WvZ6Z7niL6/K7fWTgCwBkAnAHAGQCcAcAZAJwBwD3uDzD9w+Z9/6rXyUffP6p7fUH070euB9CfDADOAOAMAM4A4AwAzgDgHtcDdO/zV33fXz0vkV0PUL3fwGi2v4DrAbTEAOAMAM4A4AwAzgDgDABuez1At+51/dnXm+neR7Hq7x0B4AwAzgDgDADOAOAMAM4A4LbXA1TL7p8/e/z29/ez51efw7jKEQDOAOAMAM4A4AwAzgDgDABu+byA7vPyo49XzxPMXr/7+VFV6x0cAeAMAM4A4AwAzgDgDADOAOC2zw3sPvdu9ve3fzfQfU5g1urn7QgAZwBwBgBnAHAGAGcAcAYAtz0P0K17/UH3eQTZfQaz3B9ASwwAzgDgDADOAOAMAM4A4F47DzCKnhdwev3B7fv+XY4AcAYAZwBwBgBnAHAGAGcAcNvzAKe/z44+3v3+1c8/vb/AlyMAnAHAGQCcAcAZAJwBwBkA3PI8wOnft2dFv4/P/k6gen+B6td/4ggAZwBwBgBnAHAGAGcAcAYA93nrenWd4QgAZwBwBgBnAHAGAGcAcAYAZwBwBgBnAHAGAPdfAAAA///T224ynvSSmgAAAABJRU5ErkJggg==';
        });
    },
    isEmpty() {
      return !(this.selectedAlgoType !== null && this.selectedAlgoName !== null);
    }
  }
};
</script>

<style scoped></style>
