# simplebank
simplebank project in golang 


read commited only prevent dirty read 


set session transaction isolation level serializable ;
select @@transaction_isolation;
start transaction ;

select * from account.account where id=1;

select  * from account.account where balance>=100;

update account.account set balance= balance-10 where id=1;

rollback ;

commit;




CREATE TABLE account (
id INT AUTO_INCREMENT PRIMARY KEY,
owner VARCHAR(100) NOT NULL,
balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,
currency VARCHAR(10) NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



CREATE TABLE account (
id SERIAL PRIMARY KEY,
owner VARCHAR(100) NOT NULL,
balance NUMERIC(15,2) NOT NULL DEFAULT 0.00,
currency VARCHAR(10) NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
