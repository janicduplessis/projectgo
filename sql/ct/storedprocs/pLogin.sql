-- --------------------------------------------------------------------------------
-- Routine DDL
-- Note: comments before and after the routine body will not be stored by the server
-- --------------------------------------------------------------------------------
DELIMITER $$

CREATE DEFINER=`ct`@`%` PROCEDURE `pLogin`(username_IN varchar(56), passwordHash_IN varchar(64))
BEGIN
	SELECT UserId, UserName, FirstName, LastName, Email
	FROM user
	WHERE UserName = username_IN
		AND PasswordHash = passwordHash_IN;
END