<template>
  <div>
    <div id="my-graph" ref="myGraph" class="myGraph"></div>
    <!--    <v-btn @click="runKruskal">Run Kruskal Algorithm</v-btn>-->
    <!--    <v-btn @click="changeData">ChangeData</v-btn>-->
  </div>
</template>

<script>
import * as d3 from 'd3';
import { kruskal } from './visualizationUtil'; // 导入 Kruskal 算法函数

export default {
  data() {
    return {
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
        { source: 0, target: 1, weight: 5 },
        { source: 0, target: 5, weight: 3 },
        { source: 1, target: 2, weight: 2 },
        { source: 1, target: 7, weight: 1 },
        { source: 2, target: 3, weight: 6 },
        { source: 2, target: 9, weight: 5 },
        { source: 3, target: 4, weight: 8 },
        { source: 3, target: 8, weight: 9 },
        { source: 4, target: 5, weight: 4 },
        { source: 4, target: 6, weight: 3 },
        { source: 6, target: 7, weight: 7 },
        { source: 6, target: 8, weight: 2 },
        { source: 7, target: 9, weight: 4 },
        { source: 8, target: 9, weight: 6 }
      ],
      selectedLinks: [], // 存储 Kruskal 算法选择的边,
      distance: 70,
      width: 800,
      simulation: null,
      selectNode1: null,
      selectNode2: null,
      preNodeColor: 'steelblue',
      transInterval: 400
    };
  },
  mounted() {
    this.drawChart();
  },
  methods: {
    // addNode(id) {
    //   this.nodes.push({ id: id });
    //   // this.drawChart();
    //   const svg = d3.select(this.$refs.myGraph);
    //   svg.selectAll('*').remove();
    //   this.drawChart();
    // },
    getGraphInfo() {
      return { nodes: this.nodes, links: this.links };
    },
    drawChart() {
      this.width = Math.min(this.$refs.myGraph.offsetWidth, 550);
      const svg = d3
        .select('#my-graph')
        .append('svg')
        .attr('width', this.width)
        .attr('height', this.width)
        .style('border-radius', '15px')
        .style('background-color', '#F5F5F5');

      const limitRange = (x) => {
        return Math.max(this.width * 0.05, Math.min(x, this.width * 0.95));
      };

      const linkGroup = svg.append('g').attr('class', 'links');

      const nodeGroup = svg.append('g').attr('class', 'nodes');

      const simulation = d3
        .forceSimulation(this.nodes)
        .force(
          'link',
          d3
            .forceLink(this.links)
            .id((d) => d.id)
            .distance(this.distance)
        )
        .force('charge', d3.forceManyBody().strength(-350))
        .force('center', d3.forceCenter(this.width / 2, this.width / 2))
        .alphaDecay(0.02)
        .alphaTarget(0.1);

      // this.simulation = simulation;

      const link = linkGroup
        .selectAll('line')
        .data(this.links)
        .enter()
        .append('line')
        .attr('stroke', 'silver')
        .attr('stroke-width', 3)
        .attr('x1', (d) => d.source.x)
        .attr('y1', (d) => d.source.y)
        .attr('x2', (d) => d.target.x)
        .attr('y2', (d) => d.target.y)
        .attr('data-weight', (d) => d.weight);

      const linkText = linkGroup
        .selectAll('text')
        .data(this.links)
        .enter()
        .append('text')
        .attr('x', (d) => (d.source.x + d.target.x) / 2)
        .attr('y', (d) => (d.source.y + d.target.y) / 2)
        .attr('text-anchor', 'middle')
        .attr('dominant-baseline', 'central')
        .attr('fill', 'black')
        .text((d) => d.weight)
        // .style('pointer-events', 'none')
        .style('user-select', 'none');

      // linkText.attr('transform', function (d) {
      //   // 计算边的中心坐标
      //   const x = (d.source.x + d.target.x) / 2;
      //   const y = (d.source.y + d.target.y) / 2;
      //   // 计算文本元素的偏移量
      //   // const angle = Math.atan2(d.target.y - d.source.y, d.target.x - d.source.x);
      //   // const offset = 13; // 文本距离边的距离
      //   // const dy = offset * Math.cos(angle);
      //   // const dx = offset * Math.sin(angle);
      //
      //   return `translate(${newx},${newy})`;
      // });

      linkGroup
        .selectAll('line')
        .on('mouseover', function (event, d) {
          d3.select(this).transition().duration(50).attr('stroke-width', 6); // 将顶点半径变大为20
          linkText
            .filter((t) => t.source === d.source && t.target === d.target)
            .transition()
            .duration(50)
            .style('font-size', '24px');
        })
        .on('mouseout', function (event, d) {
          d3.select(this).transition().duration(50).attr('stroke-width', 3); // 将顶点半径变大为20
          linkText
            .filter((t) => t.source === d.source && t.target === d.target)
            .transition()
            .duration(50)
            .style('font-size', '16px');
        });

      linkGroup
        .selectAll('text')
        .on('mouseover', function (event, d) {
          d3.select(this).transition().duration(50).style('font-size', '24px'); // 将顶点半径变大为20
          link
            .filter((t) => t.source === d.source && t.target === d.target)
            .transition()
            .duration(50)
            .attr('stroke-width', 6);
        })
        .on('mouseout', function (event, d) {
          d3.select(this).transition().duration(50).style('font-size', '16px'); // 将顶点半径变大为20
          link
            .filter((t) => t.source === d.source && t.target === d.target)
            .transition()
            .duration(50)
            .attr('stroke-width', 3);
        });

      const node = nodeGroup
        .selectAll('circle')
        .data(this.nodes)
        .enter()
        .append('circle')
        // .style('cursor', 'move')
        .attr('r', 15)
        .attr('fill', 'steelblue')
        .attr('data-id', (d) => d.id)
        .call(
          d3
            .drag()
            .on('start', (event, d) => {
              if (!event.active) {
                simulation.alphaTarget(0.3).restart();
              }
              d.fx = d.x;
              d.fy = d.y;
            })
            .on('drag', (event, d) => {
              d.fx = event.x;
              d.fy = event.y;
            })
            .on('end', (event, d) => {
              if (!event.active) {
                simulation.alphaTarget(0);
              }
              d.fx = null;
              d.fy = null;
            })
        )
        .on('click', (event, d) => {
          if (this.selectNode1 === d.id) return;
          node
            .filter((n) => n.id === this.selectNode1)
            .transition()
            .duration(100)
            .attr('fill', this.preNodeColor);
          this.preNodeColor = node.filter((n) => n.id === d.id).attr('fill');
          node
            .filter((n) => n.id === d.id)
            .transition()
            .duration(100)
            .attr('fill', 'lightgreen');
          this.selectNode1 = d.id;
        });

      const nodeLabel = nodeGroup
        .selectAll('text')
        .data(this.nodes)
        .enter()
        .append('text')
        .text((d) => d.id)
        .attr('text-anchor', 'middle')
        .attr('dominant-baseline', 'central')
        .attr('fill', 'white')
        .style('pointer-events', 'none')
        .style('user-select', 'none');

      nodeGroup
        .selectAll('circle')
        .on('mouseover', function (event, d) {
          d3.select(this).transition().duration(100).attr('r', 19); // 当前节点半径变大为20
          nodeLabel
            .filter((labelData) => labelData.id === d.id) // 仅选择当前节点对应的label
            .transition()
            .duration(100)
            .style('font-size', '20px'); // 当前节点label字体变大为20px
        })
        .on('mouseout', function (event, d) {
          d3.select(this).transition().duration(100).attr('r', 15); // 当前节点半径变回原来大小15
          nodeLabel
            .filter((labelData) => labelData.id === d.id) // 仅选择当前节点对应的label
            .transition()
            .duration(100)
            .style('font-size', '16px'); // 当前节点label字体变回原来大小12px
        });

      simulation.on('tick', () => {
        link
          .attr('x1', (d) => limitRange(d.source.x))
          .attr('y1', (d) => limitRange(d.source.y))
          .attr('x2', (d) => limitRange(d.target.x))
          .attr('y2', (d) => limitRange(d.target.y));
        node.attr('cx', (d) => limitRange(d.x)).attr('cy', (d) => limitRange(d.y));

        linkText
          .attr('x', (d) => (limitRange(d.source.x) + limitRange(d.target.x)) / 2)
          .attr('y', (d) => (limitRange(d.source.y) + limitRange(d.target.y)) / 2);

        nodeLabel.attr('x', (d) => limitRange(d.x)).attr('y', (d) => limitRange(d.y));
      });
    },
    runKruskal() {
      const selectedLinks = kruskal(this.nodes, this.links); // 运行 Kruskal 算法
      this.selectedLinks = selectedLinks;
      // 添加 Kruskal 算法动画效果
      d3.selectAll('line')
        .attr('stroke', 'gray')
        .transition()
        .duration(1000)
        .attr('stroke', (d) => {
          if (selectedLinks.includes(d)) {
            return 'red';
          } else {
            return 'gray';
          }
        });

      const node = d3.selectAll('circle').attr('fill', 'steelblue');

      const selectedNodes = new Set();
      selectedLinks.forEach((d) => {
        selectedNodes.add(d.source);
        selectedNodes.add(d.target);
      });

      node
        .filter((d) => selectedNodes.has(d))
        .transition()
        .duration(1000)
        .attr('fill', 'red');

      d3.selectAll('text')
        .attr('fill', 'white')
        .transition()
        .duration(1000)
        .attr('fill', (d) => {
          if (selectedLinks.includes(d)) {
            return 'red';
          } else {
            return 'white';
          }
        });
    },
    changeData(newNodes, newLinks) {
      this.nodes = JSON.parse(JSON.stringify(newNodes));
      this.links = JSON.parse(JSON.stringify(newLinks));
      if (this.nodes.length > 15 && this.nodes.length < 30) {
        this.distance = 70 - 2 * this.nodes.length;
      } else if (this.nodes.length > 30) {
        this.distance = 20;
      }
      const svg = d3.select(this.$refs.myGraph);
      svg.selectAll('*').remove();
      this.drawChart();
    },
    setColor(nodeColors, linkColors, nodeLabelColors, linkTextColors) {
      const nodes = d3.selectAll('circle');
      const links = d3.selectAll('line');
      const linkTexts = d3.selectAll('.link text');
      const nodeLabels = d3.selectAll('.node text');

      nodes
        .transition()
        .duration(this.transInterval)
        .attr('fill', (d, i) => (nodeColors[i] === '' ? d : nodeColors[i]));
      links
        .transition()
        .duration(this.transInterval)
        .attr('stroke', (d, i) => (linkColors[i] === '' ? d : linkColors[i]));
      linkTexts
        .transition()
        .duration(this.transInterval)
        .attr('fill', (d, i) => (linkTextColors[i] === '' ? d : linkTextColors[i]));
      nodeLabels
        .transition()
        .duration(this.transInterval)
        .attr('fill', (d, i) => (nodeLabelColors[i] === '' ? d : nodeLabelColors[i]));
    },
    resetColor() {
      const nodes = d3.selectAll('.nodes circle');
      const links = d3.selectAll('.links line');
      const linkTexts = d3.selectAll('.links text');
      const nodeLabels = d3.selectAll('.nodes text');

      nodes
        .transition()
        .duration(this.transInterval)
        .attr('fill', () => 'steelblue');
      links
        .transition()
        .duration(this.transInterval)
        .attr('stroke', () => 'gray');
      linkTexts
        .transition()
        .duration(this.transInterval)
        .attr('fill', () => 'black');
      nodeLabels
        .transition()
        .duration(this.transInterval)
        .attr('fill', () => 'white');
    }
  }
};
</script>

<style scoped>
/*.svg {*/
/*  background-color: white;*/
/*  border: 1px solid black;*/
/*}*/

/*.circle {*/
/*  cursor: move;*/
/*}*/
.myGraph {
  display: flex;
  /*align-items: flex-end;*/
  /*height: 200px;*/
  max-height: 550px;
  max-width: 550px;
}

/*.text {*/
/*  pointer-events: none;*/
/*  user-select: none;*/
/*}*/
</style>
