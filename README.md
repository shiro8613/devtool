# devtool

## 使用用途
- 開発の際に利用するスクリプトの実行
- スクリプトの並列実行（サーバーとフロントサーバーなど）

## 使用方法
- `devtool init` プロジェクトを初期化
- 完成した`devtool.yml`ファイルを編集します。
- `devtool run 名前`で実行します。

## devtool.ymlについて

### 1. sync利用 (同期処理で上から順に実行します)

例:
```yaml
scripts:
    test:
        type: sync
        command:
            echo_aaa: "echo aaa"
```

`scripts`は固定です。

それ以降はmapになっているので、以下の型で定義します。

```yaml
test:
  type: sync
  command:
    echo_aaa: "echo aaa"
```

`test:`の部分は実行の際のコマンドです。

`type:`は実行タイプを指定します。同期実行の場合は`sync` 非同期の場合は`async`です。

`command:`にはコマンドをマップで定義します。

- `echo_aaa: "echo aaa"`の場合、`echo_aaa`の部分は実行結果表示の際に横に表示される名前です（例: `[echo_aaa] aaa`）　`echo aaa`は実行するコマンドです。

### 2. async利用 （非同期実行で全部の処理を並行して実行します）

   上記のsyncの書き方と同じですが、`type:`の場所を`async`に変更します。


## インストール方法
　後で書きます
