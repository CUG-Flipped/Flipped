//
//  STInputBar.h
//  STEmojiKeyboard
//
//  Created by zhenlintie on 15/5/29.
//  Copyright (c) 2015年 sTeven. All rights reserved.
//

#import <UIKit/UIKit.h>

@interface STInputBar : UIView

+ (instancetype)inputBar;

@property (strong, nonatomic) UIButton *keyboardTypeButton;    // 按钮类型
@property (strong, nonatomic) UITextView *textView;            // 输入框
@property (strong, nonatomic) UIButton *sendButton;            // 发送按钮
@property (strong, nonatomic) UILabel *placeHolderLabel;       // 提示文字

@property (strong, nonatomic) void (^sendDidClickedHandler)(NSString *string);    // 输入的文字通过blockc获取


@property (assign, nonatomic) BOOL fitWhenKeyboardShowOrHide;

- (void)setDidSendClicked:(void(^)(NSString *text))handler;

@property (copy, nonatomic) NSString *placeHolder;

@end
