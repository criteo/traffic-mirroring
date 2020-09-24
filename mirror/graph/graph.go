package graph

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/emicklei/dot"
	"github.com/shimmerglass/http-mirror-pipeline/mirror"
)

func FromModule(m mirror.Module) *dot.Graph {
	g := dot.NewGraph(dot.Directed)
	g.Attr("nodesep", "2")
	processNode(g, nil, m)
	return g
}

func GraphSVG(g *dot.Graph, w io.Writer) error {
	fmt.Println(g.String())
	cmd := exec.Command("dot", "-Tsvg")
	cmd.Stdin = strings.NewReader(g.String())
	cmd.Stdout = w

	return cmd.Run()
}

func processNode(g *dot.Graph, parents []dot.Node, m mirror.Module) []dot.Node {
	ctx := m.Context()
	current := g.Node(ctx.Name)
	role := ctx.Role()

	if role != "virtual" {
		current.Attr("label", dot.HTML(
			fmt.Sprintf(`
				%s<BR />
				<FONT point-size="11"><B>Throughput:</B>&nbsp;&nbsp;&nbsp;%d/s</FONT><BR />
				<FONT point-size="10">%s</FONT>
			`, ctx.Name, ctx.RPS, ctx.Type),
		))
	}

	current.Attr("width", "3")

	switch ctx.Role() {
	case "virtual":
		current.Attr("shape", "underline")
	case "source":
		current.Attr("shape", "invhouse")
		current.Attr("color", "#2D232F")
		current.Attr("fillcolor", "#634B66")
		current.Attr("fontcolor", "#FFFFFF")
		current.Attr("style", "filled")
	case "sink":
		current.Attr("shape", "invhouse")
		current.Attr("color", "#144667")
		current.Attr("fillcolor", "#2176AE")
		current.Attr("fontcolor", "#FFFFFF")
		current.Attr("style", "filled")
	default:
		current.Attr("shape", "hexagon")
		current.Attr("color", "#94A0B3")
		current.Attr("fillcolor", "#D2D7DF")
		current.Attr("fontcolor", "#000000")
		current.Attr("style", "filled")
	}

	for _, p := range parents {
		g.Edge(p, current)
	}

	children := m.Children()
	leafs := []dot.Node{}

	for _, group := range children {
		parents := []dot.Node{current}
		for i, child := range group {
			parents = processNode(g, parents, child)

			if i == len(group)-1 {
				leafs = append(leafs, parents...)
			}
		}
	}

	if len(leafs) == 0 {
		leafs = append(leafs, current)
	}

	return leafs
}
