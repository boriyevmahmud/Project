select u.id,name,firstname,
lastname,bio,phonenumbers,
createdat,updateat,deletedat,
status from users u inner join 
adress a on u.id=a.users_id;