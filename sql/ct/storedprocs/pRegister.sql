-- --------------------------------------------------------------------------------
-- Routine DDL
-- Note: comments before and after the routine body will not be stored by the server
-- --------------------------------------------------------------------------------
DELIMITER $$

CREATE DEFINER=`ct`@`%` PROCEDURE `pRegister`(username_IN varchar(56), passwordHash_IN varchar(64), firstName_IN varchar(56), lastName_IN varchar(56), email_IN varchar(56))
BEGIN
INSERT 
	INTO user (UserName, PasswordHash, FirstName, LastName, Email)
	VALUES (username_IN, passwordHash_IN, firstName_IN, lastName_IN, email_IN);

END