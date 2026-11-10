# butlerjames-hello-mcp

A sample hosted integration for ButlerJames: the smallest thing you can submit under Settings → Connections → Your own integrations.

It offers one tool, `hello`, which fetches `https://example.com/` and says hello with the page's title.

## How it works

- `butlerjames.json` declares the integration: its name, start command, tools, and the only sites it may reach (`https://example.com`, GET only). It takes no key.
- `main.go` reads one JSON value per line on stdin and answers on stdout. Asked to list, it returns its tools, exactly as declared.
- The integration has no network of its own. When `hello` is called it asks for one HTTP request (`need`), ButlerJames makes it if the declaration allows it, and the response comes back on the next line.
- It has no dependencies (standard library only), so `vendor/modules.txt` is empty. ButlerJames builds it with no internet access as `/app/integration`.

To submit it, give the repo (`owner/name`), a full 40-character commit on the default branch, and "Written in: Go".

## Licence

MIT. See [LICENSE](LICENSE).
