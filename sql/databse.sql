Create Database keys;

Create Table projects (
    id bigint primary key auto_increment,
    project_name varchar(128) not null defualt '',
    register_time datetime default current_timestamp,
    create_time datetime default current_timestamp,
    update_time datetime default current_timestamp on update current_timestamp,
    uk_project_name unique (project_name)
);

create table project_tokens (
    id bigint primary key auto_increment,
    project_id bigint not null default 0,
    project_token varchar(256) not null default '',
    project_token_expire_time datetime not null default current_timestamp,
    create_time datetime default current_timestamp,
    update_time datetime default current_timestamp on update current_timestamp,
    uk_project_id_token unique (project_id, project_token),
    idx_project_id index (protject_id),
    idx_project_id_token index (project_id, project_token)
);

Create table keys (
    id bigint primary key auto_increment,
    project_id bigint not null default 0,
    project_key varchar(128) not null default '',
    project_content text not null default '',
    create_time datetime default current_timestamp,
    update_time datetime default current_timestamp on update current_timestamp,
    uk_project_key unique (project_id, project_key),
    idx_project_id index (protject_id)
);