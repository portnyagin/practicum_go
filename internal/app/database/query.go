package database

const InsertURL = "with first_insert as (\n" +
	"    insert into urls(id, correlation_id, original_url,short_url) \n" +
	"   values( nextval('seq_urls'), $2,$3,$4) \n" +
	"   RETURNING id\n)" +
	" insert into user_urls   ( url_id, user_id) \n" +
	" select id, $1 from first_insert "

const GetURLsByUserID = "select id, user_id, original_url, short_url  \n" +
	"from urls t1, user_urls t2 \n" +
	"where \nt1.id  = t2.url_id \n" +
	"and t2.user_id=$1"

const GetOriginalURLByShort = "select original_url\n " +
	"from urls t1, user_urls t2\n " +
	"where \n " +
	"t1.id  = t2.url_id\n " +
	" and t2.is_deleted=0 " +
	"and t1.short_url =$1"
const GetOriginalURLByShortForUser = "select original_url\n " +
	"from urls t1, user_urls t2\n " +
	"where \n " +
	"t1.id  = t2.url_id\n " +
	" and t2.is_deleted=0 " +
	"and t2.user_id=$1\n " +
	"and t1.short_url =$2"

const DeleteUserURL = "update user_urls t1\n" +
	"set is_deleted = 1\n" +
	"from urls t2\n" +
	"where t1.url_id  = t2.id \n" +
	"and t1.user_id = $1 " +
	"and t2.short_url  = $2 \n"
