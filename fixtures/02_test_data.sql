USE ansible

-- Create groups
INSERT INTO `group` (`id`,`name`,`variables`,`enabled`,`monitored`) VALUES (1,"group1","{\"group_var1\": {\"group_sub_var1\": \"group_sub_val1\"}}",1,1);
INSERT INTO `group` (`id`,`name`,`variables`,`enabled`,`monitored`) VALUES (2,"group2","{\"group_var2\": \"group_val2\"}",1,1);
INSERT INTO `group` (`id`,`name`,`variables`,`enabled`,`monitored`) VALUES (3,"group3","{\"group_var3\": \"group_val3\"}",1,1);
INSERT INTO `group` (`id`,`name`,`variables`,`enabled`,`monitored`) VALUES (4,"group4","{\"group_var4\": \"group_val4\"}",1,1);
INSERT INTO `group` (`id`,`name`,`variables`,`enabled`,`monitored`) VALUES (5,"group5","{\"group_var5\": \"group_val5\"}",1,1);

-- Create hosts
INSERT INTO `host` (`id`,`host`,`hostname`,`domain`,`variables`,`enabled`,`monitored`) VALUES (1,"1.1.1.1","host1","domain.local","{\"host_var1\": {\"host_sub_var1\": \"host_sub_val1\"}}",1,1);
INSERT INTO `host` (`id`,`host`,`hostname`,`domain`,`variables`,`enabled`,`monitored`) VALUES (2,"2.2.2.2","host2","domain.local","{\"host_var2\": \"host_val2\"}",1,1);
INSERT INTO `host` (`id`,`host`,`hostname`,`domain`,`variables`,`enabled`,`monitored`) VALUES (3,"3.3.3.3","host3","domain.local","{\"host_var3\": \"host_val3\"}",1,1);

-- Create host-groups
INSERT INTO `hostgroups` (`id`,`host_id`,`group_id`) VALUES (1,1,1);
INSERT INTO `hostgroups` (`id`,`host_id`,`group_id`) VALUES (2,1,4);
INSERT INTO `hostgroups` (`id`,`host_id`,`group_id`) VALUES (3,2,2);
INSERT INTO `hostgroups` (`id`,`host_id`,`group_id`) VALUES (4,3,3);

-- Create child-groups
INSERT INTO `childgroups` (`id`,`child_id`,`parent_id`) VALUES (1,3,4);
INSERT INTO `childgroups` (`id`,`child_id`,`parent_id`) VALUES (2,4,5);
