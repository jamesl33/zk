package lister

import (
	"context"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/jamesl33/zk/internal/matcher"
	"github.com/stretchr/testify/require"
)

const content = `---
title: "Note %d"
tags: ["tag%d"]
---

%s`

func BenchmarkListerVariyingNotes(b *testing.B) {
	for _, notes := range []int{10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("%d-notes", notes), func(b *testing.B) {
			tmp := b.TempDir()

			for i := 0; i < notes; i++ {
				var (
					path    = filepath.Join(tmp, fmt.Sprintf("%d.md", i))
					content = fmt.Sprintf(content, i, i%10, strconv.Itoa(i))
				)

				err := os.WriteFile(path, []byte(content), 0o644)
				require.NoError(b, err)
			}

			l, err := NewLister(WithPath(tmp))
			require.NoError(b, err)

			b.ResetTimer()

			for b.Loop() {
				iter, stop := iter.Pull2(l.Many(context.Background()))
				defer stop()

				count := 0

				for {
					n, err, ok := iter()
					if !ok {
						break
					}

					require.NoError(b, err)
					require.NotNil(b, n)
					count++
				}

				require.Equal(b, notes, count)
			}
		})
	}
}

func BenchmarkListerVaryingSize(b *testing.B) {
	const notes = 250

	for _, size := range []int{1024, 1024 * 100, 1024 * 1024, 10 * 1024 * 1024, 20 * 1024 * 1024} {
		b.Run(fmt.Sprintf("%d-bytes", size), func(b *testing.B) {
			var (
				tmp  = b.TempDir()
				body = strings.Repeat("a", size)
			)

			for i := 0; i < notes; i++ {
				var (
					path    = filepath.Join(tmp, fmt.Sprintf("%d.md", i))
					content = fmt.Sprintf(content, i, i%10, body)
				)

				err := os.WriteFile(path, []byte(content), 0o644)
				require.NoError(b, err)
			}

			l, err := NewLister(WithPath(tmp))
			require.NoError(b, err)

			b.ResetTimer()

			for b.Loop() {
				iter, stop := iter.Pull2(l.Many(context.Background()))
				defer stop()

				count := 0

				for {
					n, err, ok := iter()
					if !ok {
						break
					}

					require.NoError(b, err)
					require.NotNil(b, n)

					count++
				}

				require.Equal(b, notes, count)
			}
		})
	}
}

func BenchmarkListerVaryingSizeWithABodyMatcher(b *testing.B) {
	const notes = 250

	for _, size := range []int{1024, 1024 * 100, 1024 * 1024, 10 * 1024 * 1024, 20 * 1024 * 1024} {
		b.Run(fmt.Sprintf("%d-bytes", size), func(b *testing.B) {
			var (
				tmp  = b.TempDir()
				body = strings.Repeat("a", size)
			)

			for i := 0; i < notes; i++ {
				var (
					path    = filepath.Join(tmp, fmt.Sprintf("%d.md", i))
					content = fmt.Sprintf(content, i, i%10, body)
				)

				err := os.WriteFile(path, []byte(content), 0o644)
				require.NoError(b, err)
			}

			matcher, err := matcher.Body("", "", ".*")
			require.NoError(b, err)

			l, err := NewLister(WithPath(tmp), WithMatcher(matcher))
			require.NoError(b, err)

			b.ResetTimer()

			for b.Loop() {
				iter, stop := iter.Pull2(l.Many(context.Background()))
				defer stop()

				count := 0

				for {
					n, err, ok := iter()
					if !ok {
						break
					}

					require.NoError(b, err)
					require.NotNil(b, n)

					count++
				}

				require.Equal(b, notes, count)
			}
		})
	}
}
