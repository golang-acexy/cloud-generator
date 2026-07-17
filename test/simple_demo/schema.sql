DROP TABLE IF EXISTS demo_employee;
DROP TABLE IF EXISTS demo_department;

CREATE TABLE demo_department (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    code VARCHAR(32) NOT NULL COMMENT '部门编码',
    name VARCHAR(64) NOT NULL COMMENT '部门名称',
    description VARCHAR(255) NOT NULL DEFAULT '' COMMENT '部门描述',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，0停用',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_demo_department_code (code),
    KEY idx_demo_department_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='框架演示部门表';

CREATE TABLE demo_employee (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    department_id BIGINT NOT NULL COMMENT '所属部门ID',
    user_id BIGINT NOT NULL COMMENT '数据所属用户ID',
    employee_no VARCHAR(32) NOT NULL COMMENT '员工编号',
    name VARCHAR(64) NOT NULL COMMENT '员工姓名',
    email VARCHAR(128) NOT NULL DEFAULT '' COMMENT '邮箱',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1在职，0离职',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_demo_employee_no (employee_no),
    KEY idx_demo_employee_department_id (department_id),
    KEY idx_demo_employee_user_id (user_id),
    KEY idx_demo_employee_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='框架演示员工表';
