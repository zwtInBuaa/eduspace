function find(parent, i) {
  if (parent[i] === i) {
    return i;
  } else {
    return find(parent, parent[i]);
  }
}

function union(parent, rank, x, y) {
  const xRoot = find(parent, x);
  const yRoot = find(parent, y);

  if (rank[xRoot] < rank[yRoot]) {
    parent[xRoot] = yRoot;
  } else if (rank[xRoot] > rank[yRoot]) {
    parent[yRoot] = xRoot;
  } else {
    parent[yRoot] = xRoot;
    rank[xRoot] += 1;
  }
}

export function kruskal(nodes, links) {
  const sortedLinks = links.slice().sort((a, b) => a.value - b.value);

  const parent = {};
  const rank = {};

  nodes.forEach((node) => {
    parent[node.id] = node.id;
    rank[node.id] = 0;
  });

  const selectedLinks = [];

  sortedLinks.forEach((link) => {
    const sourceParent = find(parent, link.source.id);
    const targetParent = find(parent, link.target.id);

    if (sourceParent !== targetParent) {
      selectedLinks.push(link);
      union(parent, rank, sourceParent, targetParent);
    }
  });

  return selectedLinks;
}
