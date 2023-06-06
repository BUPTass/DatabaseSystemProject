//表格js
function creatUserTable(name,data){
  layui.use(['table', 'dropdown'], function () {
    var table = layui.table;
    var dropdown = layui.dropdown;
    // 创建渲染实例
    table.render({
      elem: name,
      data:data
    //   , url: './data/demo1.json' // 此处为静态模拟数据，实际使用时需换成真实接口
    //   , toolbar: '#toolbarDemo'
      , defaultToolbar: ['filter', 'exports', 'print', {
        title: '提示'
        , layEvent: 'LAYTABLE_TIPS'
        , icon: 'layui-icon-tips'
      }]
      , height: 'full' // 最大高度减去其他容器已占有的高度差
      , css: [ // 重设当前表格样式
        '.layui-table-tool-temp{padding-right: 145px;}'
      ].join('')
      , cellMinWidth: 80
      , totalRow: true // 开启合计行
      , page: true
      , limit: 10
      , cols:
      [[
        {
          field: "username",
          title: '账号',
          sort: true
        },
        {
          field: "level",
          title: '级别',
          templet:function(d){
            return d.level==0?'管理员':"普通用户"
          },
          sort:true
        },{
          field: "confirmed",
          title: '是否确认',
          templet:function(d){
            return d.confirmed==true?'已确认':"未确认"
          },
          sort: true
        },
        {fixed: 'right', title:'操作', width: 134, minWidth: 125, toolbar: '#barDemo'}
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
    // table.on('toolbar(test)', function (obj) {
    //   var id = obj.config.id;
    //   var checkStatus = table.checkStatus(id);
    //   var othis = lay(this);
    //   switch (obj.event) {
    //     case "manage_user":
    //       layer.open({
    //           title: '添加用户',
    //           type: 2,
    //           area: ['80%', '80%'],
    //           content: 'add_user.html'
    //         });
    //         break;
    //     case "reload_data":
    //       table.reload('test', {
    //         where: {
    //           abc: '123456'
    //         }
    //       });
    //       layer.msg('数据已刷新');
    //       break;
            
    //   };
    // });
    table.on('tool(test)', function (obj) { // 双击 toolDouble
      var data = obj.data; // 获得当前行数据
      // console.log(obj)
      if (obj.event === 'add') {

        console.log(data.username)
        var url = "http://8.130.47.185/api/add/user?username=" + data.username
        console.log(url)
        $.post(url, function (datas) {
          console.log(datas)
          layer.msg(datas)
        })
      } else if (obj.event === 'delete') {
        // 更多 - 下拉菜单
        console.log(data.username)
        var url = "http://8.130.47.185/api/delete/user?username=" + data.username
        console.log(url)
        $.ajax({
            url:url,
            dataType: "json",
            type:"delete",
            withCredentials: true,
            success:
                function (data) {
                    console.log(data)
                    layer.alert("删除用户成功！")
                }
        });
        

        // $.delete(url, function (datas) {
        //   console.log(datas)
        //   layer.msg(datas)
        // })
      }
    });
  });
}
function creatNormalTable(id,col,data)
{
    layui.use(['table', 'dropdown'], function () {
    var table = layui.table;
    var dropdown = layui.dropdown;
    // 创建渲染实例
    table.render({
      elem: id
    //   , url: './data/demo.json' // 此处为静态模拟数据，实际使用时需换成真实接口
        ,data:data
      , toolbar: '#toolbarDemo'
      , defaultToolbar: ['filter', 'exports', 'print', {
        title: '提示'
        , layEvent: 'LAYTABLE_TIPS'
        , icon: 'layui-icon-tips'
      }]
      , height: 'full' // 最大高度减去其他容器已占有的高度差
      , css: [ // 重设当前表格样式
        '.layui-table-tool-temp{padding-right: 145px;}'
      ].join('')
      , cellMinWidth: 80
      , totalRow: true // 开启合计行
      , page: true
      , limit: 10
      , cols:
      [col],
      parseData:function(d){
        console.log(d)
        return{
          "code":0,
          "msg":"",
          "count":0,
          "data":d
        }
      }
      , 
      done: function () {
        var id = this.id;
      }
      , error: function (res, msg) {
        console.log(res, msg)
      }
    });
  });
}

