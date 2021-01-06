/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/6 19:14
 **/
package _map

type Phone struct {
	Number string `validate:"e164,required"`
}
