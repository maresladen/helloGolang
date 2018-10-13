import { Param, Component, OnClickEvent } from "rainbowui-core";
import { Util } from 'rainbow-foundation-tools';
import "../plugin/ztree/css/metroStyle/metroStyle.css";
import "../plugin/ztree/js/jquery.ztree.all-3.5";
import PropTypes from 'prop-types';

export default class Tree extends Component {

    static getCheckedNode(treeId, checked) {
        return this.getTree(treeId).getCheckedNodes(checked);
    }

    static getSelectedNode(treeId) {
        return this.getTree(treeId).getSelectedNodes();
    }

    static getTree(treeId) {
        return $.fn.zTree.getZTreeObj(treeId);
    }

    constructor(props) {
        super(props);
        this.newCount = 1;
    }

    render() {
        const { searchable, style } = this.props;
        if (Util.parseBool(searchable)) {
            return (
                <div style={style}>
                    <div style={{ position: "relative" }}>
                        {this.renderPrefixIcon()}
                        <input type="text" id={this.componentId + "_keyword"} className="empty ztreesearch" placeholder="Search" />
                        <br />
                    </div>
                    <ul id={this.componentId} className="ztree" />
                </div>
            );
        } else {
            return (<ul id={this.componentId} className="ztree" style={style} />);
        }
    }

    renderPrefixIcon() {
        return (
            <span style={{ padding: "5px 10px", position: "absolute" }}>
                <span id={this.componentId + "_prefixIcon"} className="fa fa-search" style={{ cursor: "pointer" }} />
            </span>
        );
    }

    componentDidMount() {
        this.handlerComponent();
    }

    componentDidUpdate(prevProps, prevState) {
        this.handlerComponent();
    }

    handlerComponent() {
        const { dataSource,searchOnInput, searchable } = this.props;
        if (!_.isEmpty(dataSource)) {
            for (let i = 0; i < dataSource.length; i++) {
                if (dataSource[i]['OPEN']) {
                    let tempVal = dataSource[i]['OPEN'];
                    delete dataSource[i]['OPEN'];
                    dataSource[i]['open'] = tempVal;
                }
            }
        }

        const setting = {
            callback: {
                onCheck: function (event) {
                    let nodes = getTree(this.componentId).getCheckedNodes(true);
                }
            }
        };

        $.extend(true, setting, this.getViewJson(), this.getCheckJson(), this.getDataJson(), this.getEditJson(), this.getCallBackJson());
        const tree = $.fn.zTree;
        tree.init($("#" + this.componentId), setting, dataSource);
        const treeObject = tree.getZTreeObj(this.componentId);;
        _.each(treeObject.getNodes(), (node) => {
            this.setChecked(treeObject, node);
        });
        if (Util.parseBool(searchable)) {
            if (Util.parseBool(searchOnInput)) {
                $("#" + this.componentId + "_keyword").bind("input", this.searchNode.bind(this));
            } else {
                $("#" + this.componentId + "_keyword").bind("keydown", this.searchNodeOnKeyPress.bind(this));
                $("#" + this.componentId + "_keyword").bind("blur",this.searchNode.bind(this))
            }
        }
    }

    setChecked(treeObject, node) {
        if (node.checked) {
            treeObject.checkNode(node, true, true);
        }
        if (!_.isEmpty(node["children"])) {
            _.each(node["children"], (item) => {
                this.setChecked(treeObject, item);
            });
        }
    }

    getViewJson() {
        let view = {

            autoCancelSelected: true,
            dblClickExpand: false,
            expandSpeed: "fast",
            fontCss: this.getFontCSS,
            nameIsHTML: false,
            selectedMulti: false,
            showIcon: true,
            showLine: true,
            showTitle: true,
            txtSelectedEnable: false,
        }

        if (Util.parseBool(this.props.addable)) {
            view.addHoverDom = this.addHoverDom.bind(this);
            view.removeHoverDom = this.removeHoverDom.bind(this);
        } else {
            view.addHoverDom = null;
            view.removeHoverDom = null;
        }
        if (Util.parseBool(this.props.customerDom)) {
            view.addDiyDom = this.addDiyDom.bind(this);
        } else {
            view.addDiyDom = null;
        }

        return { view };
    }

    getCheckJson() {
        return {
            check: {
                autoCheckTrigger: false,
                chkboxType: { "Y": "ps", "N": "ps" },
                chkStyle: "checkbox",
                enable: Util.parseBool(this.props.checkable),
                nocheckInherit: false,
                chkDisabledInherit: false,
                radioType: "level"
            }
        };
    }

    getDataJson() {
        const { idKey, pIdKey, name, checked } = this.props;
        return {
            data: {
                keep: {
                    leaf: false,
                    parent: false
                },
                key: {
                    checked: checked,
                    children: "children",
                    name: name,
                    title: "",
                    url: "url"
                },
                simpleData: {
                    enable: true,
                    idKey: idKey,
                    pIdKey: pIdKey,
                    rootPId: null
                }
            }
        };
    }

    getEditJson() {
        const { renameTitle, removeTitle } = this.props;
        return {
            edit: {
                drag: {
                    autoExpandTrigger: true,
                    isCopy: false,
                    isMove: Util.parseBool(this.props.moveable),
                    prev: true,
                    next: true,
                    inner: true,
                    borderMax: 10,
                    borderMin: -5,
                    minMoveSize: 5,
                    maxShowNodeNum: 5,
                    autoOpenTime: 500
                },

                enable: Util.parseBool(this.props.editable),
                editNameSelectAll: true,
                removeTitle: removeTitle,
                renameTitle: renameTitle,
                showRemoveBtn: Util.parseBool(this.props.removeable),
                showRenameBtn: Util.parseBool(this.props.renameable)
            }
        };
    }

    getCallBackJson() {
        return {
            callback: {
                beforeRemove: this.beforeRemove.bind(this),//点击删除时触发，用来提示用户是否确定删除
                beforeEditName: this.beforeEditName.bind(this),//点击编辑时触发，用来判断该节点是否能编辑
                //beforeRename:this.props.beforeRename,//编辑结束时触发，用来验证输入的数据是否符合要求
                //onRemove:this.onRemove.bind(this),//删除节点后触发，用户后台操作
                //onRename:this.onRename.bind(this),//编辑后触发，用于操作后台
                //beforeDrag:beforeDrag,//用户禁止拖动节点
                onDrag: this.onDrag.bind(this),
                onDrop: this.onDrop.bind(this),
                onDragMove: this.onDragMove.bind(this),
                onNodeCreated: this.onNodeCreated.bind(this),
                onClick: this.onClickNode.bind(this),
                onCheck: this.onCheckNode.bind(this),
                beforeCheck: this.beforeCheckNode.bind(this),
                onExpand: this.props.onExpandCallback,
                onCollapse: this.props.onCollapseCallback
            }
        };
    }

    addDiyDom(treeId, treeNode) {
        var aObj = $("#" + treeNode.tId + "_a");
        if ($("#diyBtn_" + treeNode.id).length > 0) return;
        let customerDomName = this.props.customerDomName;
        var editStr = "<span id='diyBtn_space_" + treeNode.id + "' > </span>"
            + "<button type='button' class='diyBtn1' id='diyBtn_" + treeNode.id
            + "' title='" + treeNode.name + "' onfocus='this.blur();'>"+customerDomName+"</button>";
        aObj.append(editStr);
        var btn = $("#diyBtn_" + treeNode.id);
        if (btn) 
            btn.bind("click",  this.props.customerDomEvent);
    };

    addHoverDom(treeId, treeNode) {

        let _self = this;
        const { setItemCallback } = this.props;
        let sObj = $("#" + treeNode.tId + "_span");
        if (treeNode.editNameFlag || $("#addBtn_" + treeNode.tId).length > 0) {
            return;
        }
        //let addStr = "<span class='button add' id='addBtn_" + treeNode.tId + "' title='add node'></span>";
        let addStr = "<span class='button add' id='addBtn_" + treeNode.tId + "' title='add node' onfocus='this.blur();'></span>";

        sObj.after(addStr);
        let btn = $("#addBtn_" + treeNode.tId);
        if (btn) {
            btn.bind("click", function () {
                let zTree = $.fn.zTree.getZTreeObj(_self.componentId);

                if (setItemCallback) {
                    let itemParam = setItemCallback(treeNode);
                    zTree.addNodes(treeNode, {
                        id: (100 + _self.newCount++),
                        pId: treeNode.id,
                        //name: "new node" + (_self.newCount++),
                        name: itemParam.name,
                        icon: itemParam.url
                    });
                } else {

                    zTree.addNodes(treeNode, {
                        id: (100 + _self.newCount),
                        pId: treeNode.id,
                        name: "new node" + (_self.newCount++),
                    });

                }

                return false;
            });
        }
    }

    removeHoverDom(treeId, treeNode) {
        $("#addBtn_" + treeNode.tId).unbind().remove();
    }

    onCheckNode(event, treeId, treeNode) {
        event.preventDefault();
        if (this.props.onCheck != undefined) {
            this.props.onCheck(new OnClickEvent(this, event, Param.getParameter(this)), treeId, treeNode);
        }
    }

    beforeCheckNode(treeId, treeNode) {
        if (this.props.beforeCheck != undefined) {
            this.props.beforeCheck(treeId, treeNode)
        }
    }

    onNodeCreated(event, treeId, treeNode) {
        const { onNodeCreated } = this.props;
        event.preventDefault();
        if (onNodeCreated != undefined) {
            onNodeCreated(new OnClickEvent(this, event, Param.getParameter(this)), treeId, treeNode);
        }
    }

    onDragMove(event, treeId, treeNode) {
        const { onDragMove } = this.props;
        event.preventDefault();
        if (onDragMove != undefined) {
            onDragMove(new OnClickEvent(this, event, Param.getParameter(this)), treeId, treeNode);
        }
    }

    onDrag(event, treeId, treeNode) {
        const { onDrag } = this.props;
        event.preventDefault();
        if (onDrag != undefined) {
            onDrag(new OnClickEvent(this, event, Param.getParameter(this)), treeId, treeNode);
        }
    }

    onDrop(event, treeId, treeNode) {
        const { onDrop } = this.props;
        event.preventDefault();
        if (onDrop != undefined) {
            onDrop(new OnClickEvent(this, event, Param.getParameter(this)), treeId, treeNode);
        }
    }

    onRemove(e, treeId, treeNode) {
        if (treeNode.isParent) {
            let childNodes = zTree.removeChildNodes(treeNode);
            let paramsArray = new Array();
            for (let i = 0; i < childNodes.length; i++) {
                paramsArray.push(childNodes[i].id);
            }
            alert("删除父节点的id为：" + treeNode.id + "\r\n他的孩子节点有：" + paramsArray.join(","));
            return;
        }
        alert("你点击要删除的节点的名称为：" + treeNode.name + "\r\n" + "节点id为：" + treeNode.id);
    }

    beforeEditName(treeId, treeNode) {
        if (Util.parseBool(treeNode.chkDisabled)) {
            return false;
        }
        const { beforeEditName } = this.props;
        event.preventDefault();
        if (beforeEditName != undefined) {
            return beforeEditName(new OnClickEvent(this, event, Param.getParameter(this)), treeId, treeNode);
        }
    }

    beforeRemove(treeId, treeNode) {
        // return confirm("你确定要删除吗？");
        const { messageHelper, beforeRemove } = this.props;
        if (Util.parseBool(treeNode.chkDisabled)) {
            messageHelper ? (messageHelper)() : null;
            return false;
        }

        if (beforeRemove != undefined) {
            return (beforeRemove(treeId, treeNode));
        }
    }

    // beforeRename(treeId, treeNode, newName, isCancel){
    //     if(newName.length < 3){
    //         alert("名称不能少于3个字符！");
    //         return false;
    //     }
    //     return true;
    // }

    onRename(event, treeId, treeNode, isCancel) {
        alert("修改节点的id为：" + treeNode.id + "\n修改后的名称为：" + treeNode.name);
    }

    onClickNode(event, treeId, treeNode) {
        event.preventDefault();
        const { onClick } = this.props;
        if (onClick != undefined) {
            onClick(new OnClickEvent(this, event, Param.getParameter(this)), treeId, treeNode);
            //this.props.onClick(event, treeId, treeNode);
        }
    }

    searchNodeOnKeyPress(event){
        if (event.key === 'Enter') {
            this.searchNode(event);
        }
    }

    searchNode(event) {
        const treeObj = $.fn.zTree.getZTreeObj(this.componentId);
        const value = $.trim($("#" + this.componentId + "_keyword").val());

        this.updateNodes(false, treeObj.transformToArray(treeObj.getNodes()));
        if (value === "") return;

        let keyType = "name";
        var cNodes = treeObj.getNodesByParamFuzzy(keyType, value)
        this.updateNodes(true, cNodes);
        //add by guf 展开父节点
        for (var index = 0; index < cNodes.length; index++) {
            treeObj.expandNode(cNodes[index].getParentNode(), true, false, false, false);
        }
    }

    updateNodes(highlight, nodeList) {
        const treeObj = $.fn.zTree.getZTreeObj(this.componentId);
        for (let i = 0, l = nodeList.length; i < l; i++) {
            nodeList[i].highlight = highlight;
            treeObj.updateNode(nodeList[i]);
        }
    }

    getFontCSS(treeId, treeNode) {
        //return (!!treeNode.highlight) ? {color:"#A60000", "font-weight":"bold"} : {color:"#333", "font-weight":"normal"};
        return treeNode.font ? treeNode.font :
            (!!treeNode.highlight) ? { color: "#A60000", "font-weight": "bold" } : { color: "#333", "font-weight": "normal" };
    }

};


/**
 * Tree component prop types
 */
Tree.propTypes = {
    id: PropTypes.string,
    value: PropTypes.object,
    dataSource: PropTypes.object,
    style: PropTypes.string,
    checkable: PropTypes.oneOfType([PropTypes.bool, PropTypes.string]),

    searchable: PropTypes.oneOfType([PropTypes.bool, PropTypes.string]),
    searchOnInput: PropTypes.oneOfType([PropTypes.bool, PropTypes.string]),
    editable: PropTypes.oneOfType([PropTypes.bool, PropTypes.string]),
    removeable: PropTypes.oneOfType([PropTypes.bool, PropTypes.string]),
    renameable: PropTypes.oneOfType([PropTypes.bool, PropTypes.string]),
    moveable: PropTypes.oneOfType([PropTypes.bool, PropTypes.string]),
    addable: PropTypes.oneOfType([PropTypes.bool, PropTypes.string]),
    customerDom: PropTypes.oneOfType([PropTypes.bool, PropTypes.string]),

    beforeEditName: PropTypes.func,
    beforeRemove: PropTypes.func,
    onClick: PropTypes.func,
    onDrag: PropTypes.func,
    onNodeCreated: PropTypes.func,
    customerDomEvent: PropTypes.func,
    setItemCallback: PropTypes.func,
    onExpandCallback: PropTypes.func,
    onCollapseCallback: PropTypes.func,
    messageHelper: PropTypes.func,
    pIdKey: PropTypes.string,
    idKey: PropTypes.string,
    name: PropTypes.string,
    renameTitle: PropTypes.string,
    removeTitle: PropTypes.string,
    customerDom: PropTypes.string,
    checked: PropTypes.bool
};

/**
 * Get tree component default props
 */
Tree.defaultProps = {
    style: {},
    checkable: false,
    searchable: false,
    searchOnInput: false,
    editable: false,
    removeable: true,
    renameable: true,
    moveable: false,
    addable: false,
    beforeEditName: function () { },
    beforeRemove: function () { },
    onDrag: function () { },
    onNodeCreated: function () { },
    onExpandCallback: function () { },
    onCollapseCallback: function () { },
    messageHelper: function () { },
    customerDomEvent: null,
    idKey: 'id',
    pIdKey: 'pId',
    name: 'name',
    checked: true,
    customerDom: false,
    customerDomName: 'customerDom',
    renameTitle: 'rename',
    removeTitle: 'remove'
};