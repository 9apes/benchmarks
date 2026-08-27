package com.example.cachebench;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@SpringBootApplication
public class CacheBenchApplication {

    public static void main(String[] args) {
        SpringApplication.run(CacheBenchApplication.class, args);
    }

    @RestController
    static class HelloController {
        @GetMapping("/")
        String hello() {
            return "cache-bench-maven";
        }
    }
}
