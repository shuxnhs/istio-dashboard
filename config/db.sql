SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE istio_dashboard default charset utf8mb4 COLLATE utf8mb4_general_ci;

USE istio_dashboard;

-- ----------------------------
-- Table structure for kube_config
-- ----------------------------
DROP TABLE IF EXISTS `kube_config`;
CREATE TABLE `kube_config` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `cid` varchar(255) NOT NULL COMMENT '集群id',
  `description` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '描述',
  `k8s_host` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT 'k8s的apiHost地址',
  `k8s_auth_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否使用TLS验证 0-inCluster 1-basic 2-tls 3-token',
  `k8s_auth_basic` text CHARACTER SET utf8 NOT NULL COMMENT 'basic认证',
  `k8s_auth_token` text CHARACTER SET utf8 NOT NULL COMMENT 'token认证',
  `k8s_cluster_auth_data` text CHARACTER SET utf8 NOT NULL COMMENT 'tls认证',
  `k8s_client_certificate_data` text CHARACTER SET utf8 NOT NULL COMMENT 'tls认证',
  `k8s_client_key_data` text CHARACTER SET utf8 NOT NULL COMMENT 'tls认证',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0:未定义，1:正常可用，2:不可用，3:软删除',
  `kiali_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '/api/v1/namespaces/istio-system/services/kiali:http/proxy/kiali/api' COMMENT 'kiali请求地址',
  `jaeger_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '/api/v1/namespaces/istio-system/services/tracing:http-query/proxy/jaeger/api' COMMENT 'jaeger请求地址',
  `create_time` int(11) NOT NULL,
  `update_time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
