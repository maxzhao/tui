package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ZenNotes/tui/internal/backend"
)

func cmdGraph(ctx context.Context, b backend.Backend, args Args, kind string) error {
	r, ok := b.(*backend.Remote)
	if !ok {
		return errors.New("zn " + kind + " requires --server")
	}
	if kind == "orphans" {
		out, err := r.Client().Orphans(ctx)
		if err != nil {
			return err
		}
		if args.Bool("json") {
			emitJSON(out)
		} else {
			for _, p := range out {
				emitLine(p["path"])
			}
		}
		return nil
	}
	path := args.Str("path")
	if path == "" {
		path = args.Positional(0)
	}
	out, err := r.Client().References(ctx, kind, path)
	if err != nil {
		return err
	}
	if args.Bool("json") {
		emitJSON(out)
	} else {
		for _, ref := range out {
			dest := ""
			if ref.Path != nil {
				dest = *ref.Path
			}
			emitLine(fmt.Sprintf("%s  %s  %s  %s", ref.Source, ref.Kind, ref.Status, dest))
		}
	}
	return nil
}
func cmdLinks(ctx context.Context, b backend.Backend, args Args) error {
	return cmdGraph(ctx, b, args, "outgoing")
}
func cmdUnresolved(ctx context.Context, b backend.Backend, args Args) error {
	return cmdGraph(ctx, b, args, "unresolved")
}
func cmdOrphans(ctx context.Context, b backend.Backend, args Args) error {
	return cmdGraph(ctx, b, args, "orphans")
}
func cmdCommentPending(ctx context.Context, b backend.Backend, args Args) error {
	r, ok := b.(*backend.Remote)
	if !ok {
		return errors.New("zn comment pending requires --server")
	}
	author, ok := args.String("author")
	if !ok || author == "" {
		return errors.New("--author is required")
	}
	timeout := 0
	if raw, ok := args.String("timeout"); ok {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 0 || n > 86400 {
			return errors.New("--timeout must be 0..86400 seconds")
		}
		timeout = n
	}
	out, err := r.Client().PendingComments(ctx, author, args.Bool("wait"), timeout)
	if err != nil {
		return err
	}
	if args.Bool("json") {
		emitJSON(out)
	} else {
		for _, t := range out {
			emitLine(fmt.Sprintf("%s  %s  line %d", t.Path, t.ThreadID, t.AnchorLine))
		}
	}
	return nil
}
