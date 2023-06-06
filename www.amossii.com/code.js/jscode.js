
layui.use(function () {
  var upload = layui.upload;
  var element = layui.element;
  var $ = layui.$;
  // 制作多文件上传表格
  var uploadListIns = upload.render({
    elem: '#ID-upload-demo-files',
    elemList: $('#ID-upload-demo-files-list'), // 列表元素对象
    url: 'https://httpbin.org/post', // 此处用的是第三方的 http 请求演示，实际使用时改成您自己的上传接口即可。
    accept: 'file',
    multiple: true,
    number: 3,
    auto: false,
    bindAction: '#ID-upload-demo-files-action',
    choose: function (obj) {
      var that = this;
      var files = this.files = obj.pushFile(); // 将每次选择的文件追加到文件队列
      // 读取本地文件
      obj.preview(function (index, file, result) {
        var tr = $(['<tr id="upload-' + index + '">',
        '<td>' + file.name + '</td>',
        '<td>' + (file.size / 1024).toFixed(1) + 'kb</td>',
        '<td><div class="layui-progress" lay-filter="progress-demo-' + index + '"><div class="layui-progress-bar" lay-percent=""></div></div></td>',
          '<td>',
          '<button class="layui-btn layui-btn-xs demo-reload layui-hide">重传</button>',
          '<button class="layui-btn layui-btn-xs layui-btn-danger demo-delete">删除</button>',
          '</td>',
          '</tr>'].join(''));

        // 单个重传
        tr.find('.demo-reload').on('click', function () {
          obj.upload(index, file);
        });

        // 删除
        tr.find('.demo-delete').on('click', function () {
          delete files[index]; // 删除对应的文件
          tr.remove(); // 删除表格行
          // 清空 input file 值，以免删除后出现同名文件不可选
          uploadListIns.config.elem.next()[0].value = '';
        });

        that.elemList.append(tr);
        element.render('progress'); // 渲染新加的进度条组件
      });
    },
    done: function (res, index, upload) { // 成功的回调
      var that = this;
      // if(res.code == 0){ // 上传成功
      var tr = that.elemList.find('tr#upload-' + index)
        , tds = tr.children();
      tds.eq(3).html(''); // 清空操作
      delete this.files[index]; // 删除文件队列已经上传成功的文件
      return;
      //}
      this.error(index, upload);
    },
    allDone: function (obj) { // 多文件上传完毕后的状态回调
      console.log(obj)
    },
    error: function (index, upload) { // 错误回调
      var that = this;
      var tr = that.elemList.find('tr#upload-' + index);
      var tds = tr.children();
      // 显示重传
      tds.eq(3).find('.demo-reload').removeClass('layui-hide');
    },
    progress: function (n, elem, e, index) { // 注意：index 参数为 layui 2.6.6 新增
      element.progress('progress-demo-' + index, n + '%'); // 执行进度条。n 即为返回的进度百分比
    }
  });
});
layui.use(function () {
  var dropdown = layui.dropdown;
  // 渲染
  dropdown.render({
    elem: '.demo-dropdown-base', // 绑定元素选择器，此处指向 class 可同时绑定多个元素
    data: [{
      title: '网络配置信息导入',
      id: 100
    }, {
      title: 'KPI指标信息导入',
      id: 101
    }, {
      title: 'PRB干扰信息导入',
      id: 102
    }, {
      title: 'MRO数据导入',
      id: 103
    }],
    click: function (obj) {
      this.elem.find('span').text(obj.title);
    }
  });
});


//JS
layui.use(['element', 'layer', 'util'], function () {
  var element = layui.element;
  var layer = layui.layer;
  var util = layui.util;
  var $ = layui.$;
  //头部事件
  util.event('lay-header-event', {
    menuLeft: function (othis) { // 左侧菜单事件
      layer.msg('展开左侧菜单的操作', { icon: 0 });
    },
    menuRight: function () {  // 右侧菜单事件
      layer.open({
        type: 1
        , title: '更多'
        , content: '<div style="padding: 15px;">处理右侧面板的操作</div>'
        , area: ['260px', '100%']
        , offset: 'rt' //右上角
        , anim: 'slideLeft'
        , shadeClose: true
        , scrollbar: false
      });
    }
  });
});

//表格js
layui.use(['table', 'dropdown'], function () {
  var table = layui.table;
  var dropdown = layui.dropdown;
  // 创建渲染实例
  table.render({
    elem: '#test'
    , url: 'demo1.json' // 此处为静态模拟数据，实际使用时需换成真实接口
    , toolbar: '#toolbarDemo'
    , defaultToolbar: ['filter', 'exports', 'print', {
      title: '提示'
      , layEvent: 'LAYTABLE_TIPS'
      , icon: 'layui-icon-tips'
    }]
    , height: 'full-150' // 最大高度减去其他容器已占有的高度差
    , css: [ // 重设当前表格样式
      '.layui-table-tool-temp{padding-right: 145px;}'
    ].join('')
    , cellMinWidth: 80
    , totalRow: true // 开启合计行
    , page: true
    , limit: 10
    , cols: [[
      {
        field: "user_id",
        title: 'ID',
        sort: true
      },
      {
        field: "account",
        title: '账号',
        sort: false
      },
      {
        field: "password",
        title: '密码',
        sort: false
      },
      {
        field: "identity",
        title: '权限',
        sort: false
      }
      ,
    ]]
    ,
    done: function () {
      var id = this.id;
    }
    , error: function (res, msg) {
      console.log(res, msg)
    }
  });
  // 工具栏事件
  table.on('toolbar(test)', function (obj) {
    var id = obj.config.id;
    var checkStatus = table.checkStatus(id);
    var othis = lay(this);
    switch (obj.event) {
      case "manage_user":
        layer.open({
          title: '添加用户',
          type: 2,
          area: ['80%', '80%'],
          content: 'add_user.html'
        });
        break;
      case "reload_data":
        table.reload('test', {
          where: {
            abc: '123456'
          }
        });
        layer.msg('数据已刷新');
        break;

    };
  });
});
function tourl(field) {

  var keys = Object.keys(field);
  var url = ""
  for (var i = 0; i < keys.length; i++) {
    let key = keys[i];

    console.log(key, field[key]);
    url += i != 0 ? '&' : '?'
    url += key
    url += '='
    url += field[key]
  }
  // console.log("hello")
  return url;
}