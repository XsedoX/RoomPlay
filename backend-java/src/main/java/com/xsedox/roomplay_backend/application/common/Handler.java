package com.xsedox.roomplay_backend.application.common;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

import org.springframework.stereotype.Component;

/**
 * HandlerAnnotation
 */
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.CLASS)
@Component
public @interface Handler {
}
