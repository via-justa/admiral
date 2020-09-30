CREATE DATABASE IF NOT EXISTS ansible;

USE ansible;

CREATE TABLE `group` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `variables` varchar(8192) NOT NULL DEFAULT '{}',
  `enabled` tinyint(1) NOT NULL DEFAULT '0',
  `monitored` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `group_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `childgroups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `child_id` int(11) NOT NULL,
  `parent_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `childid` (`child_id`,`parent_id`),
  KEY `childgroups_child_id` (`child_id`),
  KEY `childgroups_parent_id` (`parent_id`),
  CONSTRAINT `childgroups_ibfk_2` FOREIGN KEY (`parent_id`) REFERENCES `group` (`id`),
  CONSTRAINT `childgroups_ibfk_3` FOREIGN KEY (`child_id`) REFERENCES `group` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `host` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `host` varchar(255) NOT NULL,
  `hostname` varchar(255) NOT NULL,
  `domain` varchar(250) DEFAULT NULL,
  `variables` varchar(8192) NOT NULL DEFAULT '{}',
  `enabled` tinyint(1) NOT NULL,
  `monitored` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `host_host` (`host`),
  UNIQUE KEY `host_hostname` (`hostname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `hostgroups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `host_id` int(11) NOT NULL,
  `group_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `host_id` (`host_id`,`group_id`),
  KEY `hostgroups_host_id` (`host_id`),
  KEY `hostgroups_group_id` (`group_id`),
  CONSTRAINT `hostgroups_ibfk_1` FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE,
  CONSTRAINT `hostgroups_ibfk_2` FOREIGN KEY (`group_id`) REFERENCES `group` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE OR REPLACE
ALGORITHM = UNDEFINED VIEW `inventory` AS
select
    `group`.`name` AS `group`,
    `host`.`hostname` AS `hostname`,
    ifnull(concat(`host`.`hostname`, '.', `host`.`domain`), `host`.`host`) AS `host`,
    `host`.`variables` AS `host_vars`
from
    (`group`
left join (`host`
left join `hostgroups` on
    ((`host`.`id` = `hostgroups`.`host_id`))) on
    ((`hostgroups`.`group_id` = `group`.`id`)))
where
    ((`host`.`enabled` = 1)
    and (`group`.`enabled` = 1))
order by
    `host`.`hostname`;

CREATE OR REPLACE
ALGORITHM = UNDEFINED VIEW `children` AS
select
    `childgroups`.`id` AS `relationship_id`,
    `gparent`.`name` AS `parent`,
    `gparent`.`id` AS `parent_id`,
    `gchild`.`name` AS `child`,
    `gchild`.`id` AS `child_id`
from
    (((`childgroups`
left join `group` `gparent` on
    ((`childgroups`.`parent_id` = `gparent`.`id`)))
left join `group` `gchild` on
    ((`childgroups`.`child_id` = `gchild`.`id`)))
left join `inventory` on
    ((`gchild`.`name` = `inventory`.`group`)))
where
    ((`gparent`.`enabled` = 1)
    and (`gchild`.`enabled` = 1)
    and (`inventory`.`hostname` is not null))
group by
    `gparent`.`name`,
    `gchild`.`name`
order by
    `gparent`.`name`;

