# 国旗娘Wiki自動化スクリプト
[国旗娘Wiki](https://wikiwiki.jp/kokkimusume)のページ更新を自動化するスクリプトです。

GitHub Actionsで定期的に実行されますが、Actions > Update Pages > Run workflowから手動で実行することもできます。

[chara.csv](https://github.com/taichi765/kokkimusume-wiki-automation/blob/master/data/charas.csv)をもとにページを更新しているので、国旗娘を追加したい場合や情報を更新したい場合は`chara.csv`を編集してください。

ログインしている場合は[ファイルのページ](https://github.com/taichi765/kokkimusume-wiki-automation/blob/master/data/charas.csv)の右上の方に鉛筆マークがあるので、そのボタンを押して編集モードに入った後編集して右上の`Commit Changes`と書いてある緑のボタンを押すといろいろガイドが出てきます。
ガイドに従ってプルリクエストを出してください。

## 各ディレクトリについて
- deletion-detector/
- discord-bot/
- page-updater/
- terraform/
- types/
- wikiwiki/