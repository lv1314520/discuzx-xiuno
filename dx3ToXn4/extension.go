package dx3ToXn4

//Linux 平台下设置两论坛的根目录，移动附件、头像及版块图片
//forum:修正版主 UID,icon
//group:修正可删除用户的组 id, 将XiunoBBS 将creditsfrom为0，creditsto不为0的组ID改为101，并将 user 为此组的 gid 改为101
//post:图片数及附件数从 attach 表中提取
//thread:图片数及附件数从 attach 表中提取, 修正最后发帖者及最后帖子
//user:转换完 threads 和 posts 更新统计, 修正头像avatarstatus
