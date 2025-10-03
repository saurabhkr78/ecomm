## Go E-comm Platform full stack application
1.product requirement
2.system design and architecture
3.Tech and infra
4.deployment


We are here to facilitate the buy,sale new and old/used products in certain price range. where seller can post products and advertise it for sale and buyer can search and view product as per choice and buy it online. Application should provide features to talk to buyer with seller without mediator adn buy the products by paying the price. Payment system will collect the money of the product and certain % of the final product price will be on hold for platform and rest of will be released to seller. the communication of every process will be notify through sms and email notification

use cases:
- user can buy products online 
- user can become seller
- seller can manage products and orders
-seller can collect payments and view transactions


--------          ---------          -------
| seller  | ------>| Platform| ---> | buYer  |
---------         ----------        ---------

System Design:
1. Functional Requirements
2.Non-functional requirements


1.Functional Requirements

1. Users
- user signup/login fxnality
-user verification with otp/sms
-user can become seller/buyer

2.catelog
-product listing
-manage products(crud)
-stock management

3.payment
-buyer can purchase the product using online payments(card/online banking,upi)
4.Notification
-Email 
-SMS

2.Non-Functional Requirements
1. System should be highly available in cloud with multiple region because this is a c2c portal
2.system should maintain best practices to scale horizotally at any level
3.System should design the way it can be break down to microservices 
4.loosely coupled services and communications
5.it should have mechanism for logging and monitoring to inspect services health and availability.
6.system should design with documentation for better scope of usability to understand the achitecure and business logic of the API usages.
7.Should follow CQRS



Note:
go mod init <project-name>
configs: keeps all the config files
internals: we keep whole applications and logics
infra: infra files
pkg: all the packages

since c follows c type architecture so we need to keep a main file


To run application in hot reloading mode
make server

Question:
when to follow ORM and When to follow native sql?