layui.use(['element', 'layedit', 'form', 'layer'], function(){
    var $ = layui.$;
    var element = layui.element;
    var layedit = layui.layedit;
    var form = layui.form;
    var layer = layui.layer;
    var editorIndex = null;
    form.on('submit(install)', function(data){
        let postData = data.field;
        postData.port = Number(postData.port)
        $.post("/install", postData, function (res) {
                console.log(postData);
                if(res.code === 0) {
                        layer.alert(res.msg, function(){
                                window.location.href = "/";
                        });
                } else {
                        layer.msg(res.msg);
                }
        }, 'json');
        return false;
    });

    form.on('submit(login-submit)',function(obj){            
            $.post('/admin/login', obj.field, function(res){
                    
                    if (res.code === 0){
                            layer.msg('登录成功', {
                                    offset: '15px',
                                    icon:1,
                                    time:1000
                            }, function(){
                                    window.location.href="/";
                            })
                    } else {
                            layer.msg(res.msg);
                    }
            })
    });

    if($("#text-editor").length){
            editorIndex = layedit.build('text-editor', {
                    height: 450,
                    uploadImage:{
                            url: "/attachment/upload",
                            type:"post"
                    }
            });
    }
    form.on('submit(article-publish)', function(data){
            let postData = data.field;
            postData.id = Number(postData.id);
            if(!postData.title){
                    return layer.msg("请填写文章标题");
            }

            layedit.sync(editorIndex);
            postData.content = $('#text-editor').val();
            $.post('/article/publish', postData, function(res){
                    if (res.code === 0){
                        //     window.location.href = "/article/" + res.data.id;
                        window.location.href = "/";
                    } else {
                            layer.msg(res.msg);
                    }
            }, 'json');
            return false;
    })
});

