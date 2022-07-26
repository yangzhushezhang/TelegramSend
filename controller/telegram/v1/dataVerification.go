/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package v1

type AddTelegramVerification struct {
	RobotId string `form:"robot_id" binding:"required"`
	Remark  string `form:"remark" binding:"required"`
	Token   string `form:"token" binding:"required"`
	White  string  `form:"white" binding:"required"`
}
