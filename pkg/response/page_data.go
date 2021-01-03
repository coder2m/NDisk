/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 14:40
 */
package R

type PageData struct {
	//当前页码
	PageNo int
	//每页大小
	PageSize int
	//一共的页数
	TotalPage int
	//总条数
	TotalCount int64
	//是否是第一页
	FirstPage bool
	//是否是最后一页
	LastPage bool
	//数据
	List interface{}
}

//总条数  当前页码  每页大小   数据list
func Page(count int64, pageNo int, pageSize int, list interface{}) PageData {
	tp := int(count) / pageSize
	if int(count)%pageSize > 0 {
		tp = int(count)/pageSize + 1
	}
	return PageData{PageNo: pageNo, PageSize: pageSize, TotalPage: tp, TotalCount: count, FirstPage: pageNo == 1, LastPage: pageNo == tp, List: list}
}
