output "latest_image_tag" {
  # 将来的にprodとdevで分けたらここも分岐する
  value = var.commit_sha256
}
