tag := "v0.0.5"

tag-push:
    git tag -a {{tag}} -m "Release {{tag}}"
    git push origin {{tag}}
