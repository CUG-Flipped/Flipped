//
//  TinderUITests.m
//  TinderUITests
//
//  Created by Layer on 2020/4/23.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <XCTest/XCTest.h>

@interface TinderUITests : XCTestCase

@end

@implementation TinderUITests

- (void)setUp {
    // Put setup code here. This method is called before the invocation of each test method in the class.

    // In UI tests it is usually best to stop immediately when a failure occurs.
    self.continueAfterFailure = NO;

    // In UI tests it’s important to set the initial state - such as interface orientation - required for your tests before they run. The setUp method is a good place to do this.
}

- (void)tearDown {
    // Put teardown code here. This method is called after the invocation of each test method in the class.
}

#pragma mark - 注册
- (void)testReginstration {
    
    XCUIApplication *app2 = [[XCUIApplication alloc] init];
    [app2 launch];
    XCUIApplication *app = app2;
    [app.textFields[@"Enter full name"] tap];
    
    XCUIElement *tKey = app2/*@START_MENU_TOKEN@*/.keys[@"T"]/*[[".keyboards.keys[@\"T\"]",".keys[@\"T\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [tKey tap];
    XCUIElement *eKey = app2/*@START_MENU_TOKEN@*/.keys[@"e"]/*[[".keyboards.keys[@\"e\"]",".keys[@\"e\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [eKey tap];
    XCUIElement *sKey = app2/*@START_MENU_TOKEN@*/.keys[@"s"]/*[[".keyboards.keys[@\"s\"]",".keys[@\"s\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [sKey tap];
    XCUIElement *tKey2 = app2/*@START_MENU_TOKEN@*/.keys[@"t"]/*[[".keyboards.keys[@\"t\"]",".keys[@\"t\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [tKey2 tap];
    [app.textFields[@"Enter email"] tap];
    XCUIElement *aKey = app2/*@START_MENU_TOKEN@*/.keys[@"A"]/*[[".keyboards.keys[@\"A\"]",".keys[@\"A\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [aKey tap];
    [app2/*@START_MENU_TOKEN@*/.keys[@"b"]/*[[".keyboards.keys[@\"b\"]",".keys[@\"b\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/ tap];
    
    XCUIElement *cKey = app2/*@START_MENU_TOKEN@*/.keys[@"c"]/*[[".keyboards.keys[@\"c\"]",".keys[@\"c\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [cKey tap];
    XCUIElement *moreKey = app2/*@START_MENU_TOKEN@*/.keys[@"more"]/*[[".keyboards",".keys[@\"numbers\"]",".keys[@\"more\"]"],[[[-1,2],[-1,1],[-1,0,1]],[[-1,2],[-1,1]]],[0]]@END_MENU_TOKEN@*/;
    [moreKey tap];
    XCUIElement *key = app2/*@START_MENU_TOKEN@*/.keys[@"@"]/*[[".keyboards.keys[@\"@\"]",".keys[@\"@\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key tap];
    XCUIElement *key2 = app2/*@START_MENU_TOKEN@*/.keys[@"1"]/*[[".keyboards.keys[@\"1\"]",".keys[@\"1\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key2 tap];
    XCUIElement *key3 = app2/*@START_MENU_TOKEN@*/.keys[@"6"]/*[[".keyboards.keys[@\"6\"]",".keys[@\"6\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key3 tap];
    
    XCUIElement *key4 = app2/*@START_MENU_TOKEN@*/.keys[@"3"]/*[[".keyboards.keys[@\"3\"]",".keys[@\"3\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key4 tap];
    XCUIElement *key5 = app2/*@START_MENU_TOKEN@*/.keys[@"."]/*[[".keyboards.keys[@\".\"]",".keys[@\".\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key5 tap];
    XCUIElement *moreKey2 = app2/*@START_MENU_TOKEN@*/.keys[@"more"]/*[[".keyboards",".keys[@\"letters\"]",".keys[@\"more\"]"],[[[-1,2],[-1,1],[-1,0,1]],[[-1,2],[-1,1]]],[0]]@END_MENU_TOKEN@*/;
    [moreKey2 tap];
    [cKey tap];
    XCUIElement *oKey = app2/*@START_MENU_TOKEN@*/.keys[@"o"]/*[[".keyboards.keys[@\"o\"]",".keys[@\"o\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [oKey tap];
    XCUIElement *mKey = app2/*@START_MENU_TOKEN@*/.keys[@"m"]/*[[".keyboards.keys[@\"m\"]",".keys[@\"m\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [mKey tap];
    [app.secureTextFields[@"Enter password"] tap];
    [moreKey tap];
    [key2 tap];
    [app2/*@START_MENU_TOKEN@*/.keys[@"2"]/*[[".keyboards.keys[@\"2\"]",".keys[@\"2\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/ tap];
    [key4 tap];
    [app2/*@START_MENU_TOKEN@*/.keys[@"4"]/*[[".keyboards.keys[@\"4\"]",".keys[@\"4\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/ tap];
    
    XCUIElement *key6 = app2/*@START_MENU_TOKEN@*/.keys[@"5"]/*[[".keyboards.keys[@\"5\"]",".keys[@\"5\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key6 tap];
    [key3 tap];
    [app2/*@START_MENU_TOKEN@*/.staticTexts[@"Register"]/*[[".buttons[@\"Register\"].staticTexts[@\"Register\"]",".staticTexts[@\"Register\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/ tap];
    
    XCUIElement *svprogresshudElement = app.otherElements[@"SVProgressHUD"];
    [svprogresshudElement tap];
    
}

#pragma mark - 登录
- (void)testLogin {
    
    XCUIApplication *app2 = [[XCUIApplication alloc] init];
    [app2 launch];
    XCUIElement *enterEmailTextField = app2.textFields[@"Enter email"];
    [enterEmailTextField tap];
    
    XCUIElement *tKey = app2/*@START_MENU_TOKEN@*/.keys[@"T"]/*[[".keyboards.keys[@\"T\"]",".keys[@\"T\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [tKey tap];
    XCUIElement *eKey = app2/*@START_MENU_TOKEN@*/.keys[@"e"]/*[[".keyboards.keys[@\"e\"]",".keys[@\"e\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [eKey tap];
    XCUIElement *sKey = app2/*@START_MENU_TOKEN@*/.keys[@"s"]/*[[".keyboards.keys[@\"s\"]",".keys[@\"s\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [sKey tap];
    XCUIElement *tKey2 = app2/*@START_MENU_TOKEN@*/.keys[@"t"]/*[[".keyboards.keys[@\"t\"]",".keys[@\"t\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [tKey2 tap];
    
    XCUIApplication *app = app2;
    [app.secureTextFields[@"Enter password"] tap];
    
    XCUIElement *moreKey = app2/*@START_MENU_TOKEN@*/.keys[@"more"]/*[[".keyboards",".keys[@\"numbers\"]",".keys[@\"more\"]"],[[[-1,2],[-1,1],[-1,0,1]],[[-1,2],[-1,1]]],[0]]@END_MENU_TOKEN@*/;
    [moreKey tap];
    
    XCUIElement *key = app2/*@START_MENU_TOKEN@*/.keys[@"1"]/*[[".keyboards.keys[@\"1\"]",".keys[@\"1\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key tap];
    [app2/*@START_MENU_TOKEN@*/.keys[@"2"]/*[[".keyboards.keys[@\"2\"]",".keys[@\"2\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/ tap];
    
    XCUIElement *key2 = app2/*@START_MENU_TOKEN@*/.keys[@"3"]/*[[".keyboards.keys[@\"3\"]",".keys[@\"3\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key2 tap];
    [app2/*@START_MENU_TOKEN@*/.keys[@"4"]/*[[".keyboards.keys[@\"4\"]",".keys[@\"4\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/ tap];
    
    XCUIElement *key3 = app2/*@START_MENU_TOKEN@*/.keys[@"5"]/*[[".keyboards.keys[@\"5\"]",".keys[@\"5\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key3 tap];
    XCUIElement *key4 = app2/*@START_MENU_TOKEN@*/.keys[@"6"]/*[[".keyboards.keys[@\"6\"]",".keys[@\"6\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key4 tap];
    XCUIElement *key5 = app2/*@START_MENU_TOKEN@*/.keys[@"7"]/*[[".keyboards.keys[@\"7\"]",".keys[@\"7\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key5 tap];
    
    XCUIElement *key6 = app2/*@START_MENU_TOKEN@*/.keys[@"8"]/*[[".keyboards.keys[@\"8\"]",".keys[@\"8\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key6 tap];
    [app2/*@START_MENU_TOKEN@*/.keys[@"9"]/*[[".keyboards.keys[@\"9\"]",".keys[@\"9\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/ tap];
    
    XCUIElement *key7 = app2/*@START_MENU_TOKEN@*/.keys[@"0"]/*[[".keyboards.keys[@\"0\"]",".keys[@\"0\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key7 tap];
    
    XCUIElement *logInButton = app.buttons[@"Log In"];
    [logInButton tap];
    [enterEmailTextField tap];
    [moreKey tap];
    
    XCUIElement *key8 = app2/*@START_MENU_TOKEN@*/.keys[@"@"]/*[[".keyboards.keys[@\"@\"]",".keys[@\"@\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key8 tap];
    [key tap];
    [key4 tap];
    [key2 tap];
    XCUIElement *key9 = app2/*@START_MENU_TOKEN@*/.keys[@"."]/*[[".keyboards.keys[@\".\"]",".keys[@\".\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [key9 tap];
    XCUIElement *moreKey2 = app2/*@START_MENU_TOKEN@*/.keys[@"more"]/*[[".keyboards",".keys[@\"letters\"]",".keys[@\"more\"]"],[[[-1,2],[-1,1],[-1,0,1]],[[-1,2],[-1,1]]],[0]]@END_MENU_TOKEN@*/;
    [moreKey2 tap];
    XCUIElement *cKey = app2/*@START_MENU_TOKEN@*/.keys[@"c"]/*[[".keyboards.keys[@\"c\"]",".keys[@\"c\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [cKey tap];
    XCUIElement *oKey = app2/*@START_MENU_TOKEN@*/.keys[@"o"]/*[[".keyboards.keys[@\"o\"]",".keys[@\"o\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [oKey tap];
    XCUIElement *mKey = app2/*@START_MENU_TOKEN@*/.keys[@"m"]/*[[".keyboards.keys[@\"m\"]",".keys[@\"m\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/;
    [mKey tap];
    [logInButton tap];
    [app2/*@START_MENU_TOKEN@*/.buttons[@"Clear text"]/*[[".textFields[@\"Enter email\"].buttons[@\"Clear text\"]",".buttons[@\"Clear text\"]"],[[[-1,1],[-1,0]]],[0]]@END_MENU_TOKEN@*/ tap];
    [moreKey tap];
    [key tap];
    [key5 tap];
    [key4 tap];
    [key2 tap];
    [key3 tap];
    [key2 tap];
    [key4 tap];
    [key4 tap];
    [key tap];
    [key5 tap];
    [key5 tap];
    [key8 tap];
    [key tap];
    [key4 tap];
    [key2 tap];
    [key9 tap];
    [moreKey2 tap];
    [cKey tap];
    [oKey tap];
    [mKey tap];
    [logInButton tap];
}

#pragma mark - 主页测试
- (void)testHome {
    XCUIApplication *app = [[XCUIApplication alloc] init];
    [app launch];
    XCUIElement *element = [[[[[[[[[[app childrenMatchingType:XCUIElementTypeWindow] elementBoundByIndex:0] childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther].element childrenMatchingType:XCUIElementTypeOther] elementBoundByIndex:1];
    [element swipeLeft];
    [element swipeRight];
    [element swipeUp];
    [element swipeDown];
    [element swipeLeft];
    [element swipeRight];
    
    XCUIElement *twoimgButton = app.buttons[@"twoImg"];
    [twoimgButton tap];
    
    XCUIElement *fourimgButton = app.buttons[@"fourImg"];
    [fourimgButton tap];
    [twoimgButton tap];
    [fourimgButton tap];
    
    XCUIElement *threeimgButton = app.buttons[@"threeImg"];
    [threeimgButton tap];
    [threeimgButton tap];
    [app.buttons[@"firstImg"] tap];
    [twoimgButton tap];
    [fourimgButton tap];
    
    // Use recording to get started writing UI tests.
    // Use XCTAssert and related functions to verify your tests produce the correct results.
}

#pragma mark - 个人中心
- (void)testinfo {
    
}

#pragma mark - 通讯录
- (void)testAddressBook {
    
}

#pragma mark - 聊天
- (void)testAddressChat {
    
}


- (void)testLaunchPerformance {
    if (@available(macOS 10.15, iOS 13.0, tvOS 13.0, *)) {
        // This measures how long it takes to launch your application.
        [self measureWithMetrics:@[XCTOSSignpostMetric.applicationLaunchMetric] block:^{
            [[[XCUIApplication alloc] init] launch];
        }];
    }
}
@end
