# Mustache Engine

<p align="center">
  <video width="320" height="240" controls>
  <source src="https://media.istockphoto.com/id/1170803160/nl/video/snor-icoon-vonken-deeltjes-op-zwarte-achtergrond.mp4?s=mp4-640x640-is&k=20&c=1ILsyAbNi9qWI4Oj3lXA-zLx-pLVPB49Xn7qkJSXGuY=">
</video>
  <strong>Mustache - A Workflow Engine for Offensive Security</strong>
  </p>
</p>

***

## 🔥 What is Mustache?

Mustache is a Workflow Engine for Offensive Security. It was designed to build a foundation with the capability and
flexibility that allows you to build your own reconnaissance system and watch tower and run it on a large number of targets.


## 📦 Installation

### Build the engine from the source

Make sure you installed `golang >= v1.17`

```bash
go install -v github.com/omidxplimbo/mustache@latest
```

## 🚀 Key Features of Osmedeus

- [x] Significantly speed up your recon process
- [x] Organize your scan results
- [x] Efficiently to customize and optimize your recon process
- [x] Seamlessly integrate with new public and private tools
- [x] Easy to scale across large number of targets
- [x] Easy to synchronize the results across many places

## 💡 Usage

```bash
# Scan Usage:
mustache -module continues -flow [flowName] -target [targetName] -path [pathFile] -db [databaseName]
mustache -module recon -flow [flowName] -target [targetName] -path [pathFile] -db [databaseName]
mustache -module project -flow init -project [projectName] 
mustache -module project -flow backup -project [projectName]
mustache -module project -flow remove -project [projectName]
mustache -module project -flow backup -project [projectName] -collection [collectionName]
```


## License

`Mustache` is made with ♥ by [@omidxplimbo] and it is released under the MIT license.
