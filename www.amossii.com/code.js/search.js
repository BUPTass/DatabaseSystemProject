function renderdropdown(url, id) {
    layui.use(function () {
        var dropdown = layui.dropdown;
        var data_list = []
        $.ajax({
            url: url,
            dataType: "json",
            withCredentials: true,
            success:
                function (data) {
                    for (var i = 0, len = data.length; i < len; i++) {
                        var name = data[i];
                        data_list.push({ title: name })
                    }
                    data_list = data
                    console.log(data_list);
                }
        });
        dropdown.render({
            elem: id,
            data: data_list,
            click: function (obj) {
                this.elem.val(obj.title);
            },
            style: 'min-width: 235px;'
        });
    });
}
function getCookie(cname) {
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i].trim();
        if (c.indexOf(name) == 0) return c.substring(name.length, c.length);
    }
    return "";
}