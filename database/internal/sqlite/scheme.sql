CREATE TABLE IF NOT EXISTS `group` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `name` varchar(255) NOT NULL,
  `variables` varchar(8192) NOT NULL DEFAULT '{}',
  `enabled` tinyint(1) NOT NULL DEFAULT '0',
  `monitored` tinyint(1) NOT NULL DEFAULT '0',
--   `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--   `updated` TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE (`name`)
);

CREATE TABLE IF NOT EXISTS `childgroups` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `child_id` integer NOT NULL,
  `parent_id` integer NOT NULL,
  UNIQUE (`child_id`,`parent_id`),
  CONSTRAINT `childgroups_ibfk_2` FOREIGN KEY (`parent_id`) REFERENCES `group` (`id`),
  CONSTRAINT `childgroups_ibfk_3` FOREIGN KEY (`child_id`) REFERENCES `group` (`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `host` (
  `id` integer  NOT NULL PRIMARY KEY AUTOINCREMENT,
  `host` varchar(255) NOT NULL,
  `hostname` varchar(255) NOT NULL,
  `domain` varchar(250) DEFAULT NULL,
  `variables` longtext NOT NULL DEFAULT '{}',
  `enabled` tinyint(1) NOT NULL,
  `monitored` tinyint(1) NOT NULL DEFAULT '0',
--   `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--   `updated` TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE (`host`),
  UNIQUE (`hostname`)
);

CREATE TABLE IF NOT EXISTS `hostgroups` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `host_id` integer NOT NULL,
  `group_id` integer NOT NULL,
  UNIQUE (`host_id`),
  FOREIGN KEY (`host_id`) REFERENCES `host` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`group_id`) REFERENCES `group` (`id`) ON DELETE CASCADE
);

CREATE VIEW IF NOT EXISTS `hostgroup_view` AS
SELECT
    `hostgroups`.`id` AS `relationship_id`,
    `host`.`hostname` AS `host`,
    `hostgroups`.`host_id` AS `host_id`,
    `group`.`name` AS `group`,
    `hostgroups`.`group_id` AS `group_id`
FROM `hostgroups`
LEFT JOIN `group`
    ON `hostgroups`.`group_id` = `group`.`id`
LEFT JOIN `host`
    ON `hostgroups`.`host_id` = `host`.`id`
ORDER BY `group`.`name`;

CREATE VIEW IF NOT EXISTS `childgroups_view` AS
SELECT
    `childgroups`.`id` AS `relationship_id`,
    `gparent`.`name` AS `parent`,
    `gparent`.`id` AS `parent_id`,
    `gchild`.`name` AS `child`,
    `gchild`.`id` AS `child_id`
FROM `childgroups`
LEFT JOIN `group` `gparent`
	ON `childgroups`.`parent_id` = `gparent`.`id`
LEFT JOIN `group` `gchild`
	ON `childgroups`.`child_id` = `gchild`.`id`
ORDER BY `gparent`.`name`;

CREATE VIEW IF NOT EXISTS `host_view` AS
WITH RECURSIVE inherited (child_id, parent_id) AS (
SELECT
	child_id,
	parent_id
from
	childgroups_view cv
UNION ALL
SELECT
	cv.child_id,
	i.parent_id
FROM
	inherited i
JOIN childgroups_view cv ON
	i.child_id = cv.parent_id )
SELECT
	`host`.`id` AS `host_id`,
    `host`.`hostname` AS `hostname`,
    `host`.`domain` AS `domain`,
    `host`.`host` AS `host`,
    `host`.`enabled` AS `enabled`,
    `host`.`monitored` AS `monitored`,
    `host`.`variables` AS `variables`,
    ifnull(group_concat(distinct `g1`.`name`),"") AS `direct_group`,
    ifnull(group_concat(distinct `g2`.`name`),"") AS `inherited_groups`
FROM
	host
LEFT JOIN hostgroup_view hv ON
	host.id = hv.host_id
LEFT JOIN inherited i ON
	hv.group_id = i.child_id
LEFT JOIN `group` g1 ON
	hv.group_id = g1.id
LEFT JOIN `group` g2 ON
	i.parent_id = g2.id
GROUP BY
	host.hostname
ORDER BY
	host.hostname;

CREATE VIEW IF NOT EXISTS `groups_view` AS 
WITH RECURSIVE inherited(`child_id`, `parent_id`) AS (
SELECT
    `cv`.`child_id` AS `child_id`,
    `cv`.`parent_id` AS `parent_id`
FROM
    `childgroups_view` `cv`
UNION ALL
SELECT
    `cv`.`child_id` AS `child_id`,
    `i`.`parent_id` AS `parent_id`
FROM
    `inherited` `i`
JOIN `childgroups_view` `cv` ON `i`.`child_id` = `cv`.`parent_id`)
SELECT 
	`g1`.`id` AS `group_id`,
    `g1`.`name` AS `name`,
    `g1`.`enabled` AS `enabled`,
    `g1`.`monitored` AS `monitored`,
	COUNT(DISTINCT h.id) AS `num_hosts`,
	COUNT(DISTINCT `g2`.`name`) AS `num_children`,
	IFNULL(GROUP_CONCAT(DISTINCT `g2`.`name`),'') AS `child_groups`,
	`g1`.`variables` AS `variables`
FROM `group` g1
LEFT JOIN `inherited` `i` ON
     `g1`.`id` = `i`.`parent_id`
LEFT JOIN `group` `g2` ON
     `i`.`child_id` = `g2`.`id`
LEFT JOIN hostgroups hg ON
	`hg`.`group_id` = `g1`.`id` 
LEFT JOIN host h ON
	`h`.`id` = `hg`.`host_id` 
GROUP BY `g1`.`id` 
ORDER BY `g1`.`id`;
