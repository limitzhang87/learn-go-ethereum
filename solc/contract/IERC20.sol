// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface IERC20 {
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    /**
    * @dev 返回总供给
    */
    function totalSupply() external view returns (uint256);

    /**
    * @dev 返回账户`account`所持有的代币数
    */
    function balanceOf(address account) external view returns(uint256);


    /**
    * @dev 转账`account` 单位代币， 从调用者账户到另一个账户
    *
    * 成功就返回true
    *
    * 释放{Transfer}事件.
    */
    function transfer(address to, uint256 value) external returns(bool);

    /**
    * @dev 返回`owner`授权给`spender`的额度
    * 当{approve}或{transferfrom}被调用时，`allowance`需要改变
    */
    function allowance(address owner, address spender) external view returns(uint256);

    /**
    * @dev 调用者账户给`spender`账户授权`amount`数量代币
    *
    * 成功返回true
    *
    * 释放{Approval}事件
    */
    function approve(address spender, uint256 amount) external returns(bool);

    function transferFrom(address from, address to, uint256 amount) external returns(bool);
}
