INSERT INTO users (id, name, email, password_hash) VALUES ('b0c88c13-a070-459a-b93a-785cbe39d113', 'User1 Surname', 'user1@example.com', '$2a$12$thbAlwUzDjDRGBU0D7OSNu9oh9Lmmu0635cUP1lnISoihli6AXOu2');
INSERT INTO users (id, name, email, password_hash) VALUES ('f8f6b175-1096-4773-88b0-d3f1be0c062d', 'User2 Surname', 'user2@example.com', '$2a$12$thbAlwUzDjDRGBU0D7OSNu9oh9Lmmu0635cUP1lnISoihli6AXOu2');
INSERT INTO users (id, name, email, password_hash) VALUES ('ccba9088-60ff-4726-b53d-45471181a30e', 'User3 Surname', 'user3@example.com', '$2a$12$thbAlwUzDjDRGBU0D7OSNu9oh9Lmmu0635cUP1lnISoihli6AXOu2');
INSERT INTO user_sessions (session_id, user_id, expires) VALUES ('session-1', 'b0c88c13-a070-459a-b93a-785cbe39d113', '0001-01-01 00:00:00');
INSERT INTO user_sessions (session_id, user_id, expires) VALUES ('session-2', 'f8f6b175-1096-4773-88b0-d3f1be0c062d', '0001-01-01 00:00:00');

INSERT INTO polls (id, name, type, options) VALUES ('b517b6d9-a25e-4b02-97a2-31fec1fd8afc', 'simple-majority-1000', 'simple-marjority', 'optA,opt-2,something else');

INSERT INTO votes (user_id, poll_id, option) VALUES ('b0c88c13-a070-459a-b93a-785cbe39d113', 'b517b6d9-a25e-4b02-97a2-31fec1fd8afc', 'opt-2');
INSERT INTO votes (user_id, poll_id, option) VALUES ('ccba9088-60ff-4726-b53d-45471181a30e', 'b517b6d9-a25e-4b02-97a2-31fec1fd8afc', 'opt-2');
INSERT INTO votes (user_id, poll_id, option) VALUES ('f8f6b175-1096-4773-88b0-d3f1be0c062d', 'b517b6d9-a25e-4b02-97a2-31fec1fd8afc', 'optA');